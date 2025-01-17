namespace server.src
{
    /// <summary>
    /// Represents a message with a correlation ID, content, sender, recipient, and timestamp.
    /// </summary>
    /// <param name="correlationId">The correlation ID of the message.</param>
    /// <param name="content">The content of the message.</param>
    /// <param name="from">The sender of the message.</param>
    /// <param name="to">The recipient of the message.</param>
    /// <param name="time">The timestamp of the message.</param>
    internal class Message(uint correlationId, string content, string from, string to, ulong time)
    {
        public uint CorrelationId { get; } = correlationId;
        public string Content { get; } = content;
        public string From { get; } = from;
        public string To { get; } = to;
        public ulong Time { get; } = time;
    }
}
