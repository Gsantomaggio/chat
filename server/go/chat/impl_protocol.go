package chat

import (
	"bufio"
	"gsantomaggio/chat/server/internal"
)

// CommandLogin is a command to login into the chat server.

type CommandLogin struct {
	correlationId uint32 // 4 bytes
	username      string // max 256 characters, for example "gabriele" [8, 103, 97, 98, 114, 105, 101, 108, 101]
}

func NewCommandLoginWithCorrelation(username string, correlationId uint32) *CommandLogin {
	return &CommandLogin{username: username, correlationId: correlationId}
}

func NewCommandLogin(username string) *CommandLogin {
	return &CommandLogin{username: username}
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
	return chatProtocolUint32 + // correlationId
		chatProtocolSizeUint16 + // size of the string
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
	From          string
	To            string
	Time          uint64
}

func NewCommandMessage(message, from string, to string, time uint64) *CommandMessage {
	return &CommandMessage{Message: message, From: from, To: to, Time: time}
}

func NewCommandMessageWithCorrelationId(message, from string, to string, correlationId uint32, time uint64) *CommandMessage {
	return &CommandMessage{Message: message, From: from, To: to, correlationId: correlationId, Time: time}
}

func (m *CommandMessage) Read(reader *bufio.Reader) error {
	return readMany(reader, &m.correlationId, &m.Message, &m.From, &m.To, &m.Time)
}

func (m *CommandMessage) Key() uint16 {
	return CommandMessageKey
}

func (m *CommandMessage) SizeNeeded() int {
	return chatProtocolUint32 + // correlationId
		chatProtocolUint64 + // time
		chatProtocolSizeUint16 + // size of the string message
		len(m.Message) + // actual size of the message
		chatProtocolSizeUint16 + // size of the string from
		len(m.From) + // actual size of the "from"
		chatProtocolSizeUint16 + // size of the string to
		len(m.To) // actual size of the "to"
}

func (m *CommandMessage) CorrelationId() uint32 {
	return m.correlationId
}

func (m *CommandMessage) SetCorrelationId(id uint32) {
	m.correlationId = id
}

func (m *CommandMessage) Version() int16 {
	return Version1
}

func (m *CommandMessage) Write(writer *bufio.Writer) (int, error) {
	return writeMany(writer, m.correlationId, m.Message, m.From, m.To, m.Time)
}

// ChatHeader is the header of the chat protocol.
type ChatHeader struct {
	// total size of this header + command content
	command uint16 // 2 bytes
	version int16  // 2 bytes
}

func NewChatHeaderFromCommand(command internal.CommandWrite) *ChatHeader {
	return &ChatHeader{command: command.Key(), version: command.Version()}
}

func NewChatHeader(version int16, command uint16) *ChatHeader {
	return &ChatHeader{command: command, version: version}
}

func (c *ChatHeader) Write(writer *bufio.Writer) (int, error) {
	return writeMany(writer, c.version, c.command)
}

func (c *ChatHeader) Read(reader *bufio.Reader) error {
	return readMany(reader, &c.version, &c.command)
}

func (c *ChatHeader) Key() uint16 {
	return c.command
}

func (c *ChatHeader) Version() int16 {
	return c.version
}

func (c *ChatHeader) SizeNeeded() int {
	return chatProtocolHeaderSizeBytes
}

type GenericResponse struct {
	correlationId uint32
	responseCode  uint16
}

func (g *GenericResponse) Key() uint16 {
	return GenericResponseKey
}

func (g *GenericResponse) SizeNeeded() int {
	return chatProtocolUint32 + // correlationId
		chatProtocolSizeUint16 // responseCode
}

func (g *GenericResponse) Version() int16 {
	return Version1
}

func (g *GenericResponse) SetCorrelationId(id uint32) {
	g.correlationId = id
}

func NewGenericResponse(responseCode uint16) *GenericResponse {
	return &GenericResponse{
		responseCode: responseCode,
	}
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
