using System.Net.Sockets;
using System.Text;

namespace server.src
{
    /// <summary>
    /// Handles incoming messages from a network stream, processes user login and message commands,
    /// and sends appropriate responses.
    /// </summary>
    internal class MessageHandler(NetworkStream stream, BinaryReader reader)
    {
        private static readonly Encoding encoding = Encoding.UTF8;
        private readonly BinaryReader _reader = reader;
        private SingleUser? _user = null;
        private readonly NetworkStream _stream = stream;

        /// <summary>
        /// Processes incoming messages and handles user login or message commands.
        /// </summary>
        /// <param name="usr">The user associated with the message, if any.</param>
        /// <returns>The updated user after processing the message.</returns>
        public async Task<SingleUser?> HandleMessage(SingleUser? usr)
        {
            _user = usr;
            bool isValidMessage = ReadMessageLength(out var length, out var remaining);
            if (!isValidMessage)
            {
                Logger.LogError("Message not correct, declared len {length}, but remaining len {remaining}", length, remaining);
                return _user;
            }

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

        /// <summary>
        /// Reads the length of the incoming message and calculates the remaining bytes in the stream.
        /// </summary>
        /// <param name="length">The length of the incoming message.</param>
        /// <param name="remaining">The number of remaining bytes in the stream.</param>
        /// <returns>True if the length of the message matches the remaining bytes in the stream; otherwise, false.</returns>
        public bool ReadMessageLength(out uint length, out long remaining)
        {
            length = _reader.ReadUInt32BE();
            var position = _reader.BaseStream.Position;
            remaining = _reader.BaseStream.Length - position;
            return length == remaining;
        }

        /// <summary>
        /// Reads the header of the incoming message to determine the command key.
        /// </summary>
        /// <returns>The command key from the message header.</returns>
        public ushort ReadHeader()
        {
            _reader.ReadByte(); // version
            return _reader.ReadUInt16BE();
        }

        /// <summary>
        /// Reads a message from the stream and constructs a Message object.
        /// </summary>
        /// <param name="id">The correlation ID of the message.</param>
        /// <returns>The constructed Message object.</returns>
        public Message ReadMessage(uint id)
        {
            string content = ReadString();
            string from = ReadString();
            string to = ReadString();
            ulong time = _reader.ReadUInt64BE();
            return new Message(id, content, from, to, time);
        }

        /// <summary>
        /// Reads a string from the stream.
        /// </summary>
        /// <returns>The string read from the stream.</returns>
        public string ReadString()
        {
            ushort len = _reader.ReadUInt16BE();
            return encoding.GetString(_reader.ReadBytes(len));
        }

        /// <summary>
        /// Sends all queued messages for a user over the network stream.
        /// </summary>
        /// <param name="stream">The network stream to send messages over.</param>
        /// <param name="user">The user whose messages are to be sent.</param>
        public static async Task SendUserMessages(NetworkStream stream, SingleUser user)
        {
            while (user.Messages.Count > 0)
            {
                var message = user.Messages.Dequeue();
                await SendMessage(stream, message);
            }
        }

        /// <summary>
        /// Sends a single message to a user. If the user is offline, queues the message.
        /// </summary>
        /// <param name="user">The recipient user.</param>
        /// <param name="message">The message to be sent.</param>
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

        /// <summary>
        /// Sends a message over the network stream.
        /// </summary>
        /// <param name="stream">The network stream to send the message over.</param>
        /// <param name="m">The message to be sent.</param>
        private static async Task SendMessage(NetworkStream stream, Message m)
        {
            byte[] messageBytes = Commands.CreateCommandMessage(m);
            await stream.WriteAsync(messageBytes);
        }

        /// <summary>
        /// Sends a response with a correlation ID and response code over the network stream.
        /// </summary>
        /// <param name="id">The correlation ID of the response.</param>
        /// <param name="code">The response code.</param>
        private async Task SendResponse(uint id, ushort code)
        {
            var responseBytes = Commands.CreateResponse(id, code);
            await _stream.WriteAsync(responseBytes);
            Logger.LogInformation("Response sent with code {code} and correlationId {id}", code, id);
        }
    }
}
