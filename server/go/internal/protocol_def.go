package internal

import (
	"bufio"
)

type CommandWrite interface {
	Write(writer *bufio.Writer) (int, error)
	Key() uint16
	// SizeNeeded must return the size required to encode this CommandWrite
	// plus the size of the Header. The size of the Header is always 4 bytes
	SizeNeeded() int
	Version() byte
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

type CommandRead interface {
	Read(reader *bufio.Reader) error
	Key() uint16
}

type ResponseWrite = SyncCommandWrite
type ResponseRead = CommandRead
