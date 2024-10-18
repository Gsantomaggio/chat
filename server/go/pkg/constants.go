package pkg

const (
	CommandLoginKey   uint16 = 0x01
	CommandMessageKey uint16 = 0x02
	Version1          int16  = 1

	chatProtocolHeaderSize = chatProtocolHeaderSizeBytes +
		chatProtocolCorrelationIdSizeBytes
	chatProtocolKeySizeBytes       = 2
	chatProtocolKeySizeUint8       = 1
	chatProtocolKeySizeUint16      = 2
	chatProtocolKeySizeInt         = 4
	chatProtocolStringLenSizeBytes = 2

	chatProtocolVersionSizeBytes       = 2
	chatProtocolCorrelationIdSizeBytes = 4
	chatProtocolKeySizeUint32          = 4

	chatProtocolHeaderSizeBytes = 4
)

// / response codes
const (
	ResponseCodeOk    uint16 = 0x0001
	ResponseCodeError uint16 = 0x0002
)
