package pkg

import (
	"bufio"
	"bytes"
	"encoding"
	"sync"
)

const (
	CommandLoginKey   uint16 = 0x00010
	CommandMessageKey uint16 = 0x00011
	Version1          int16  = 1

	chatProtocolHeaderSize = chatProtocolKeySizeBytes +
		chatProtocolVersionSizeBytes +
		chatProtocolCorrelationIdSizeBytes
	chatProtocolKeySizeBytes       = 2
	chatProtocolKeySizeUint8       = 1
	chatProtocolKeySizeUint16      = 2
	chatProtocolStringLenSizeBytes = 2

	chatProtocolVersionSizeBytes       = 2
	chatProtocolCorrelationIdSizeBytes = 4
	chatProtocolKeySizeUint32          = 4
)

type CommandWrite interface {
	Write(writer *bufio.Writer) (int, error)
	Key() uint16
	// SizeNeeded must return the size required to encode this CommandWrite
	// plus the size of the Header. The size of the Header is always 4 bytes
	SizeNeeded() int
	Version() int16
}

// SyncCommandWrite is the interface that wraps the WriteTo method.
// The interface is implemented by all commands that are sent to the server.
// and that have responses in RPC style.
// SetCorrelationId CorrelationId is used to match the response with the request.
type SyncCommandWrite interface {
	CommandWrite // Embedding the CommandWrite interface
	SetCorrelationId(id uint32)
	CorrelationId() uint32
}

// WriteCommand sends the Commands to the server.
// The commands are sent in the following order:
// 1. Command
// 2. Flush
// The flush is required to make sure that the commands are sent to the server.
// WriteCommand doesn't care about the response.
var mutex = &sync.Mutex{} // it is needed because the bufio.Writer is not thread safe
func WriteCommand[T CommandWrite](request T, writer *bufio.Writer) error {
	mutex.Lock()
	defer mutex.Unlock()

	bWritten, err := request.Write(writer)
	if err != nil {
		return err
	}
	if (bWritten) != (request.SizeNeeded()) {
		panic("WriteTo Command: Not all bytes written")
	}
	return writer.Flush()
}

type Command interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}

type CommandLogin struct {
	correlationId uint32
	username      string
}

func NewCommandLogin(username string, correlationId uint32) *CommandLogin {
	return &CommandLogin{username: username, correlationId: correlationId}
}

func (l *CommandLogin) Username() string {
	return l.username
}

func (l *CommandLogin) GetCorrelationId() uint32 {
	return l.correlationId
}

func (l *CommandLogin) UnmarshalBinary(data []byte) error {
	buff := bytes.NewReader(data)
	rd := bufio.NewReader(buff)
	return readMany(rd, &l.correlationId, &l.username)
}

func (l *CommandLogin) Key() uint16 {
	return CommandLoginKey
}

func (l *CommandLogin) SizeNeeded() int {
	return chatProtocolHeaderSize +
		chatProtocolKeySizeUint32 + // correlationId
		chatProtocolKeySizeUint16 + // size of the string
		len(l.username)
}

func (l *CommandLogin) SetCorrelationId(id uint32) {
	l.correlationId = id
}

func (l *CommandLogin) CorrelationId() uint32 {
	return l.correlationId
}

func (l *CommandLogin) Version() int16 {
	return Version1
}

func (l *CommandLogin) Write(writer *bufio.Writer) (int, error) {
	return writeMany(writer, l.correlationId, l.username)
}

type CommandMessage struct {
	correlationId uint32
	Message       string
	To            string
}

func NewCommandMessage(message, to string) *CommandMessage {
	return &CommandMessage{Message: message, To: to}
}

func (m *CommandMessage) UnmarshalBinary(data []byte) error {
	buff := bytes.NewReader(data)
	rd := bufio.NewReader(buff)
	return readMany(rd, &m.correlationId, &m.Message, &m.To)
}

func (m *CommandMessage) Key() uint16 {
	return CommandMessageKey
}

func (m *CommandMessage) SizeNeeded() int {
	return chatProtocolHeaderSize +
		chatProtocolStringLenSizeBytes +
		len(m.Message) +
		chatProtocolStringLenSizeBytes +
		len(m.To)
}

func (m *CommandMessage) SetCorrelationId(id uint32) {
	m.correlationId = id
}

func (m *CommandMessage) CorrelationId() uint32 {
	return m.correlationId
}

func (m *CommandMessage) Version() int16 {
	return Version1
}

func (m *CommandMessage) Write(writer *bufio.Writer) (int, error) {
	return writeMany(writer, m.correlationId, m.Message, m.To)
}

// ChatHeader is the header of the chat protocol.
type ChatHeader struct {
	// total size of this header + command content
	length int
	// Key ID
	command uint16
	version int16
}
