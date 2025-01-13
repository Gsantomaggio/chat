using System.Text;

namespace server.src
{
    /// <summary>
    /// Provides methods for creating command messages and responses for network communication.
    /// </summary>
    internal static class Commands
    {
        private static readonly Encoding encoding = Encoding.UTF8;

        /// <summary>
        /// Creates a command message from a Message object.
        /// </summary>
        /// <param name="message">The message to be converted into a command message.</param>
        /// <returns>A byte array representing the command message.</returns>
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

        /// <summary>
        /// Creates a response message with a correlation ID and response code.
        /// </summary>
        /// <param name="id">The correlation ID of the response.</param>
        /// <param name="code">The response code.</param>
        /// <returns>A byte array representing the response message.</returns>
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

        /// <summary>
        /// Serializes a string into a byte array with its length prefixed.
        /// </summary>
        /// <param name="value">The string to be serialized.</param>
        /// <returns>A byte array representing the serialized string.</returns>
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

        /// <summary>
        /// Calculates the total length of a message based on its components.
        /// </summary>
        /// <param name="contentLength">The length of the content string.</param>
        /// <param name="fromLength">The length of the from string.</param>
        /// <param name="toLength">The length of the to string.</param>
        /// <returns>The total length of the message.</returns>
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
                   Constants.protocolUint64SizeBytes;
        }

        /// <summary>
        /// Writes the header of a message to a memory stream.
        /// </summary>
        /// <param name="stream">The memory stream to write to.</param>
        /// <param name="messageLength">The length of the message.</param>
        /// <param name="key">The command key.</param>
        /// <param name="correlationId">The correlation ID of the message.</param>
        private static void WriteHeader(MemoryStream stream, uint messageLength, ushort key, uint correlationId)
        {
            stream.Write(EndianHelpers.GetBytesUInt32BE(messageLength));
            stream.Write([Constants.Version]);
            stream.Write(EndianHelpers.GetBytesUInt16BE(key));
            stream.Write(EndianHelpers.GetBytesUInt32BE(correlationId));
        }

        /// <summary>
        /// Writes the payload of a message to a memory stream.
        /// </summary>
        /// <param name="stream">The memory stream to write to.</param>
        /// <param name="content">The content of the message.</param>
        /// <param name="from">The sender of the message.</param>
        /// <param name="to">The recipient of the message.</param>
        /// <param name="time">The timestamp of the message.</param>
        private static void WriteMessagePayload(MemoryStream stream, byte[] content, byte[] from, byte[] to, ulong time)
        {
            stream.Write(content);
            stream.Write(from);
            stream.Write(to);
            stream.Write(EndianHelpers.GetBytesUInt64BE(time));
        }

        /// <summary>
        /// Writes the payload of a response message to a memory stream.
        /// </summary>
        /// <param name="stream">The memory stream to write to.</param>
        /// <param name="code">The response code.</param>
        private static void WriteResponsePayload(MemoryStream stream, ushort code)
        {
            stream.Write(EndianHelpers.GetBytesUInt16BE(code));
        }
    }
}