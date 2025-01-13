namespace server.src
{
    internal static class Constants
    {
        internal const ushort CommandLoginKey = 0x01;
        internal const ushort CommandMessageKey = 0x02;
        internal const ushort CommandResponseKey = 0x03;
        
        internal const byte Version = 1;

        internal const int protocolHeaderSizeBytes = protocolVersionSizeByte + protocolKeySizeBytes;
        internal const int protocolVersionSizeByte = 1;
        internal const int protocolKeySizeBytes = 2;
        internal const int protocolStringLenSizeBytes = 2;
        internal const int protocolCorrelationIdSizeBytes = 4;
        internal const int protocolUint8SizeBytes = 1;
        internal const int protocolUint16SizeBytes = 2;
        internal const int protocolUint32SizeBytes = 4;
        internal const int protocolUint64SizeBytes = 8;
                       
        internal const ushort ResponseCodeOk = 0x01;
        internal const ushort ResponseCodeErrorUserNotFound = 0x03;
        internal const ushort ResponseCodeErrorUserAlreadyLogged = 0x04;
    }

}
