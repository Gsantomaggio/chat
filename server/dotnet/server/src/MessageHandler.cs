using System.Net.Sockets;
using System.Text;
using Microsoft.Extensions.Logging;

namespace server.src
{
    internal class MessageHandler(NetworkStream stream, BinaryReader reader)
    {
        private static readonly Encoding encoding = Encoding.UTF8;
        private readonly BinaryReader _reader = reader;
        private SingleUser? _user = null;
        private readonly NetworkStream _stream = stream;

        public async Task<SingleUser?> HandleMessage(SingleUser? usr)
        {
            _user = usr;
            ReadMessageLength();
            ushort key = ReadHeader();
            uint correlationId = _reader.ReadUInt32BE();

            switch (key)
            {
                case Constants.CommandLoginKey:
                    string username = ReadString();
                    ushort responseCode = Users.Login(username, _stream, out SingleUser newUser);

                    await SendResponse(correlationId, responseCode);

                    if (responseCode == 4)
                    {
                        Logger.LogWarning("User {username} already logged in", username);
                        return _user;
                    }

                    await SendUserMessages(_stream, newUser);

                    return newUser;

                case Constants.CommandMessageKey:
                    ushort respCode = Constants.ResponseCodeErrorUserNotFound;
                    if (_user == null)
                    {
                        await SendResponse(correlationId, respCode);
                        return _user;
                    }
                    Message message = ReadMessage(correlationId);

                    if (Users.Instance.TryGetValue(message.To, out SingleUser user))
                    {
                        respCode = Constants.ResponseCodeOk;
                        await SendResponse(correlationId, respCode);
                        await SendSingleMessage(user, message);
                    }
                    else
                    {
                        Logger.LogWarning("User {receiver} not exists", message.To);
                        await SendResponse(correlationId, respCode);
                    }
                    return _user;

                default:
                    Logger.LogError("Received wrong COMMAND in the header. KEY: {key}", key);
                    throw new Exception($"Received wrong COMMAND in the header. KEY: {key}");
            }
        }

        public void ReadMessageLength()
        {
            uint length = _reader.ReadUInt32BE();
            var position = _reader.BaseStream.Position;
            var remaining = _reader.BaseStream.Length - position;
            if (length != remaining)
            {
                Logger.LogError("Message not correct, declared len {length}, but remaining len {remaining}", length, remaining);
                throw new Exception($"Message not correct, declared len {length}, but remaining len {remaining}");
            }
        }

        public ushort ReadHeader()
        {
            _reader.ReadByte(); // version
            return _reader.ReadUInt16BE();
        }

        public Message ReadMessage(uint id)
        {
            string content = ReadString();
            string from = ReadString();
            string to = ReadString();
            ulong time = _reader.ReadUInt64BE();
            return new Message(id, content, from, to, time);
        }

        public string ReadString()
        {
            ushort len = _reader.ReadUInt16BE();
            return encoding.GetString(_reader.ReadBytes(len));
        }

        public static async Task SendUserMessages(NetworkStream stream, SingleUser user)
        {
            while (user.Messages.Count > 0)
            {
                var message = user.Messages.Dequeue();
                await SendMessage(stream, message);
            }
        }

        public static async Task SendSingleMessage(SingleUser user, Message message)
        {
            if (user.Status == "online" && user.Stream != null)
            {
                await SendMessage(user.Stream, message);
            }
            else
            {
                user.Messages.Enqueue(message);
                Logger.LogWarning("User {username} is offline and received a message from {sender}", user.Username, message.From);
            }
        }

        private static async Task SendMessage(NetworkStream stream, Message m)
        {
            byte[] messageBytes = Commands.CreateCommandMessage(m);
            await stream.WriteAsync(messageBytes);
        }

        private async Task SendResponse(uint id, ushort code)
        {
            var responseBytes = Commands.CreateResponse(id, code);
            await _stream.WriteAsync(responseBytes);
            Logger.LogInformation("Response sent with code {code} and correlationId {id}", code, id);
        }
    }
}
