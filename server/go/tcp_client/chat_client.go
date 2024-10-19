package tcp_client

import (
	"bufio"
	"fmt"
	"gsantomaggio/chat/server/chat"
	"gsantomaggio/chat/server/internal"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type Response struct {
	responseCode  int
	data          chan any
	correlationId uint32
}

func NewResponse(correlationId uint32) *Response {
	return &Response{
		correlationId: correlationId,
		data:          make(chan any),
	}
}

type ChatClient struct {
	tcpConn           *net.TCPConn
	chMessages        chan *chat.CommandMessage
	nextCorrelationId uint32
	respMutex         sync.Mutex
	responses         map[uint32]*Response
	currentUser       string
}

func NewChatClient(receiver chan *chat.CommandMessage) *ChatClient {
	fc := &ChatClient{
		chMessages: receiver,
		responses:  make(map[uint32]*Response),
	}
	return fc
}
func (f *ChatClient) atomicIncrementCorrelationId() uint32 {
	return atomic.AddUint32(&f.nextCorrelationId, 1)
}

func (f *ChatClient) AddResponse(correlationId uint32) {
	f.respMutex.Lock()
	defer f.respMutex.Unlock()
	f.responses[correlationId] = NewResponse(correlationId)
}

func (f *ChatClient) GetResponse(correlationId uint32) *Response {
	f.respMutex.Lock()
	defer f.respMutex.Unlock()
	return f.responses[correlationId]
}

func (f *ChatClient) WaitResponse(correlationId uint32) (any, error) {
	resp := f.GetResponse(correlationId)
	if resp == nil {
		return nil, fmt.Errorf("Response not found for correlationId %d\n", correlationId)
	}
	select {
	case data := <-resp.data:
		return data, nil
	case <-time.After(time.Duration(5) * time.Second):
		return nil, fmt.Errorf("Timeout waiting for response with correlationId %d\n", correlationId)
	}

}

func (f *ChatClient) RemoveResponse(correlationId uint32) {
	f.respMutex.Lock()
	defer f.respMutex.Unlock()
	if _, ok := f.responses[correlationId]; !ok {
		fmt.Printf("Response not found for correlationId %d\n", correlationId)
		return
	}
	close(f.responses[correlationId].data)
	delete(f.responses, correlationId)
}

func (f *ChatClient) Connect(servAddr string) error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
	if err != nil {
		return err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return err
	}
	f.tcpConn = conn

	go func() {
		for {
			err := f.WaitMessages()
			if err != nil {
				return
			}
		}
	}()

	return nil
}

func (f *ChatClient) Close() error {
	return f.tcpConn.Close()
}

func (f *ChatClient) sendRPCCommand(command internal.SyncCommandWrite) (*chat.GenericResponse, error) {
	command.SetCorrelationId(f.atomicIncrementCorrelationId())
	f.AddResponse(command.CorrelationId())
	err := chat.WriteCommandWithHeader(command, bufio.NewWriter(f.tcpConn))
	if err != nil {
		return nil, err
	}
	resp, err := f.WaitResponse(command.CorrelationId())
	if err != nil {
		return nil, err
	}
	return resp.(*chat.GenericResponse), nil
}

func (f *ChatClient) Login(user string) (*chat.GenericResponse, error) {
	commandLogin := chat.NewCommandLogin(user)
	f.currentUser = user
	return f.sendRPCCommand(commandLogin)
}

func (f *ChatClient) SendMessage(message string, to string) (*chat.GenericResponse, error) {
	commandMessage := chat.NewCommandMessage(message, f.currentUser, to)
	return f.sendRPCCommand(commandMessage)
}

func (f *ChatClient) ReadMessage(reader *bufio.Reader) (*chat.CommandMessage, error) {
	msg := &chat.CommandMessage{}
	err := msg.Read(reader)
	return msg, err
}

func (f *ChatClient) WaitMessages() error {
	reader := bufio.NewReader(f.tcpConn)
	for {
		header := &chat.ChatHeader{}
		err := header.Read(reader)
		if err != nil {
			return err
		}
		switch header.Key() {
		case chat.CommandMessageKey:
			{
				msg, err := f.ReadMessage(reader)
				if err != nil {
					return nil
				}

				f.chMessages <- msg

			}
		case chat.GenericResponseKey:
			{
				generic := &chat.GenericResponse{}
				err := generic.Read(reader)
				if err != nil {
					return err
				}
				res := f.GetResponse(generic.CorrelationId())
				res.data <- generic
			}

		}

		return nil
	}
}
