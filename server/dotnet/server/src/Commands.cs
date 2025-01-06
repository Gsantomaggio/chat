using System.Text;

namespace server.src
{
    internal static class Commands
    {
        private static readonly Encoding encoding = Encoding.UTF8;

        public static byte[] CreateCommandMessage(Message message)
        {
            var contentBytes = SerializeString(message.Content);
            var fromBytes = SerializeString(message.From);
            var toBytes = SerializeString(message.To);

            uint totalLength = CalculateTotalMessageLength(contentBytes.Length, fromBytes.Length, toBytes.Length);

            byte[] bufferMessage = new byte[4+totalLength];
            using var stream = new MemoryStream(bufferMessage);
            WriteHeader(stream, totalLength, Constants.CommandMessageKey, message.CorrelationId);
            WriteMessagePayload(stream, contentBytes, fromBytes, toBytes, message.Time);

            var comMex = stream.ToArray();
            stream.Close();

            return comMex;
        }

        public static byte[] CreateResponse(uint id, ushort code)
        {
            const uint totalLength = 9;
            const ushort key = Constants.CommandResponseKey;

            byte[] bufferResponse = new byte[4+totalLength];
            using var stream = new MemoryStream(bufferResponse);
            WriteHeader(stream, totalLength, key, id);
            WriteResponsePayload(stream, code);

            var resp = stream.ToArray();
            stream.Close();

            return resp;
        }

        private static byte[] SerializeString(string value)
        {
            byte[] contentBytes = encoding.GetBytes(value);
            ushort length = (ushort)encoding.GetByteCount(value);
            byte[] lengthBytes = EndianHelpers.GetBytesUInt16BE(length);

            byte[] result = new byte[lengthBytes.Length + contentBytes.Length];
            Buffer.BlockCopy(lengthBytes, 0, result, 0, lengthBytes.Length);
            Buffer.BlockCopy(contentBytes, 0, result, lengthBytes.Length, contentBytes.Length);

            return result;
        }

        private static uint CalculateTotalMessageLength(int contentLength, int fromLength, int toLength)
        {
            return Constants.protocolHeaderSizeBytes +
                   Constants.protocolCorrelationIdSizeBytes +
                   Constants.protocolStringLenSizeBytes +
                   (uint)contentLength +
                   Constants.protocolStringLenSizeBytes +
                   (uint)fromLength +
                   Constants.protocolStringLenSizeBytes +
                   (uint)toLength +
                   Constants.protocolUint64;
        }

        private static void WriteHeader(MemoryStream stream, uint messageLength, ushort key, uint correlationId)
        {
            stream.Write(EndianHelpers.GetBytesUInt32BE(messageLength));
            stream.Write([Constants.Version]);
            stream.Write(EndianHelpers.GetBytesUInt16BE(key));
            stream.Write(EndianHelpers.GetBytesUInt32BE(correlationId));
        }

        private static void WriteMessagePayload(MemoryStream stream, byte[] content, byte[] from, byte[] to, ulong time)
        {
            stream.Write(content);
            stream.Write(from);
            stream.Write(to);
            stream.Write(EndianHelpers.GetBytesUInt64BE(time));
        }

        private static void WriteResponsePayload(MemoryStream stream, ushort code)
        {
            stream.Write(EndianHelpers.GetBytesUInt16BE(code));
        }
    }
}