package chat

const (
	CommandLoginKey    uint16 = 0x01
	CommandMessageKey  uint16 = 0x02
	GenericResponseKey uint16 = 0x03
	Version1           int16  = 1

	chatProtocolHeaderSizeBytes = chatProtocolVersionSizeBytes + // version
		chatProtocolKeySizeBytes // command
	chatProtocolKeySizeBytes       = 2
	chatProtocolKeySizeUint8       = 1
	chatProtocolSizeUint16         = 2
	chatProtocolKeySizeInt         = 4
	chatProtocolStringLenSizeBytes = 2

	chatProtocolVersionSizeBytes       = 2
	chatProtocolCorrelationIdSizeBytes = 4
	chatProtocolUint32                 = 4
	chatProtocolUint64                 = 8
)

// / response codes
const (
	ResponseCodeOk uint16 = 0x0001
	//ResponseCodeError                   uint16 = 0x0002
	ResponseCodeErrorUserNotFound      uint16 = 0x0003
	ResponseCodeErrorUserAlreadyLogged uint16 = 0x0004
)
