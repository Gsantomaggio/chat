package chat

import (
	"bufio"
	"gsantomaggio/chat/server/internal"
	"sync"
)

// WriteCommand sends the Commands to the server.
// The commands are sent in the following order:
// 1. Command
// 2. Flush
// The flush is required to make sure that the commands are sent to the server.
// WriteCommand doesn't care about the response.
var mutex = &sync.Mutex{} // it is needed because the bufio.Writer is not thread safe
func WriteCommand[T internal.CommandWrite](request T, writer *bufio.Writer) error {
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

func WriteCommandWithHeader[T internal.CommandWrite](request T, writer *bufio.Writer) error {
	mutex.Lock()
	defer mutex.Unlock()
	hr := NewChatHeaderFromCommand(request)
	hWritten, err := hr.Write(writer)
	if err != nil {
		return err
	}
	bWritten, err := request.Write(writer)
	if err != nil {
		return err
	}
	if (bWritten + hWritten) != (request.SizeNeeded() + 4) {
		panic("WriteTo Command: Not all bytes written")
	}
	return writer.Flush()
}

func WriteResponse[T internal.ResponseWrite](response T, writer *bufio.Writer) error {
	mutex.Lock()
	defer mutex.Unlock()

	bWritten, err := response.Write(writer)
	if err != nil {
		return err
	}
	if (bWritten) != (response.SizeNeeded()) {
		panic("WriteTo Response: Not all bytes written")
	}
	return writer.Flush()
}
