namespace server.src
{
    internal class Message(uint correlationId, string content, string from, string to, ulong time)
    {
        public uint CorrelationId { get; } = correlationId;
        public string Content { get; } = content;
        public string From { get; } = from;
        public string To { get; } = to;
        public ulong Time { get; } = time;
    }
}
