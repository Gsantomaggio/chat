namespace server.src
{
    internal static class Constants
    {
        internal const ushort CommandLoginKey = 0x01;
        internal const ushort CommandMessageKey = 0x02;
        internal const ushort CommandResponseKey = 0x03;
        
        internal const byte Version = 1;

        internal const uint protocolHeaderSizeBytes = protocolVersionSizeByte + protocolKeySizeBytes;
        internal const uint protocolKeySizeBytes = 2;
        internal const uint protocolKeySizeUint8 = 1;
        internal const uint protocolSizeUint16 = 2;
        internal const uint protocolKeySizeInt = 4;
        internal const uint protocolStringLenSizeBytes = 2;
        internal const uint protocolVersionSizeByte = 1;
        internal const uint protocolCorrelationIdSizeBytes = 4;
        internal const uint protocolUint32 = 4;
        internal const uint protocolUint64 = 8;
                       
        internal const ushort ResponseCodeOk = 0x01;
        internal const ushort ResponseCodeErrorUserNotFound = 0x03;
        internal const ushort ResponseCodeErrorUserAlreadyLogged = 0x04;
    }

}
