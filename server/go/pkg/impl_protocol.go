package pkg

import (
	"bufio"
	"gsantomaggio/chat/server/internal"
)

// CommandLogin is a command to login into the chat server.

type CommandLogin struct {
	correlationId uint32 // 4 bytes
	username      string // max 256 characters, for example "gabriele" [8, 103, 97, 98, 114, 105, 101, 108, 101]
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

func (l *CommandLogin) Key() uint16 {
	return CommandLoginKey
}

func (l *CommandLogin) SizeNeeded() int {
	return chatProtocolHeaderSize +
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

func (l *CommandLogin) Read(reader *bufio.Reader) error {
	return readMany(reader, &l.correlationId, &l.username)
}

/// ***** END LOGIN ***

type CommandMessage struct {
	correlationId uint32
	Message       string
	To            string
}

func NewCommandMessage(message, to string, correlationId uint32) *CommandMessage {
	return &CommandMessage{Message: message, To: to, correlationId: correlationId}
}

func (m *CommandMessage) Read(reader *bufio.Reader) error {
	return readMany(reader, &m.correlationId, &m.Message, &m.To)
}

func (m *CommandMessage) Key() uint16 {
	return CommandMessageKey
}

func (m *CommandMessage) SizeNeeded() int {
	return chatProtocolHeaderSize +
		chatProtocolKeySizeUint16 +
		len(m.Message) +
		chatProtocolKeySizeUint16 +
		len(m.To)
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
	length  int    // 4 bytes
	command uint16 // 2 bytes
	version int16  // 2 bytes
}

func NewChatHeaderFromCommand(command internal.CommandWrite) *ChatHeader {
	return &ChatHeader{length: command.SizeNeeded(), command: command.Key(), version: command.Version()}
}

func NewChatHeader(length int, version int16, command uint16) *ChatHeader {
	return &ChatHeader{length: length, command: command, version: version}
}

func (c *ChatHeader) Write(writer *bufio.Writer) (int, error) {
	return writeMany(writer, c.length, c.version, c.command)
}

func (c *ChatHeader) Read(reader *bufio.Reader) error {
	return readMany(reader, &c.length, &c.version, &c.command)

}

func (c *ChatHeader) Key() uint16 {
	return c.command
}

func (c *ChatHeader) Version() int16 {
	return c.version
}

func (c *ChatHeader) Length() int {
	return c.length
}

type GenericResponse struct {
	correlationId uint32
	responseCode  uint16
}

func NewGenericResponse(correlationId uint32, responseCode uint16) *GenericResponse {
	return &GenericResponse{correlationId: correlationId, responseCode: responseCode}
}

func (g *GenericResponse) CorrelationId() uint32 {
	return g.correlationId
}

func (g *GenericResponse) ResponseCode() uint16 {
	return g.responseCode
}

func (g *GenericResponse) Write(writer *bufio.Writer) (int, error) {
	return writeMany(writer, g.correlationId, g.responseCode)
}

func (g *GenericResponse) Read(reader *bufio.Reader) error {
	return readMany(reader, &g.correlationId, &g.responseCode)
}
