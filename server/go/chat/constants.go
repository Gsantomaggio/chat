package chat

const (
	CommandLoginKey    uint16 = 0x01
	CommandMessageKey  uint16 = 0x02
	GenericResponseKey uint16 = 0x03
	Version1           byte   = 1

	chatProtocolHeaderSizeBytes = chatProtocolVersionSizeByte + // version
		chatProtocolKeySizeBytes // command
	chatProtocolKeySizeBytes       = 2
	chatProtocolKeySizeUint8       = 1
	chatProtocolSizeUint16         = 2
	chatProtocolKeySizeInt         = 4
	chatProtocolStringLenSizeBytes = 2

	chatProtocolVersionSizeByte        = 1
	chatProtocolCorrelationIdSizeBytes = 4
	chatProtocolUint32                 = 4
	chatProtocolUint64                 = 8
)

// / response codes
const (
	ResponseCodeOk uint16 = 0x01
	//ResponseCodeError                   uint16 = 0x0002
	ResponseCodeErrorUserNotFound      uint16 = 0x03
	ResponseCodeErrorUserAlreadyLogged uint16 = 0x04
)
