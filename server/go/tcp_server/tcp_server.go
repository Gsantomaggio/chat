package tcp_server

import (
	"bufio"
	"errors"
	"fmt"
	"gsantomaggio/chat/server/chat"
	"io"
	"net"
)

type TcpServerer interface {
}

type TcpServer struct {
	host     string
	port     int
	users    map[string]*User
	listener net.Listener
	chEvents chan *Event
}

func NewTcpServer(host string, port int, events chan *Event) *TcpServer {
	return &TcpServer{
		host:     host,
		port:     port,
		users:    make(map[string]*User),
		chEvents: events,
	}
}

func (t *TcpServer) DispatchEvent(message string, isAnError bool, level int) {
	if t.chEvents != nil {
		t.chEvents <- NewEvent(message, isAnError, level)
	}
}

func (t *TcpServer) StartInAThread() error {
	go func() {
		err := t.Start()
		if err != nil {
			t.DispatchEvent(fmt.Sprintf("Error starting server: %v", err), true, 1)
		}
	}()
	return nil
}
func (t *TcpServer) Start() error {
	address := fmt.Sprintf("%s:%d", t.host, t.port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		t.DispatchEvent(fmt.Sprintf("Error starting server: %v", err), true, 2)
		return fmt.Errorf("error starting TCP server: %v", err)
	}
	t.listener = listener

	t.DispatchEvent(fmt.Sprintf("Server started at %s", address), false, 2)
	for {
		conn, err := listener.Accept()
		if err != nil {
			t.DispatchEvent(fmt.Sprintf("Error accepting connection: %v", err), true, 3)
			break
		}
		go t.handleConnection(conn)
	}

	t.DispatchEvent("Server stopped", false, 2)
	return nil
}

func (t *TcpServer) Stop() error {
	return t.listener.Close()
}

func (t *TcpServer) handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	var user *User
	for {

		readerFull, err := chat.ReadFullBufferFromSource(reader)
		if err != nil {
			t.DispatchEvent(fmt.Sprintf("Error reading source: %v", err), true, 3)
			return
		}

		header := &chat.ChatHeader{}
		err = header.Read(readerFull)
		if err != nil {
			if errors.Is(err, io.EOF) {
				t.DispatchEvent("Connection closed due of EOF", false, 2)
			} else {
				t.DispatchEvent(fmt.Sprintf("Error reading header: %v", err), true, 3)
			}
			break
		}
		var correlationId uint32
		var lastSendError error
		switch header.Key() {
		case chat.CommandLoginKey:
			login := &chat.CommandLogin{}
			err := login.Read(readerFull)
			if err != nil {
				t.DispatchEvent(fmt.Sprintf("Error reading login: %v", err), true, 3)
				break
			}
			correlationId = login.CorrelationId()
			t.DispatchEvent(fmt.Sprintf("Login request for user %s", login.Username()), false, 1)
			if t.users[login.Username()] != nil && t.users[login.Username()].IsOnLine() {
				t.DispatchEvent(fmt.Sprintf("User %s already logged", login.Username()), false, 4)
				lastSendError = t.sendBackResponse(chat.ResponseCodeErrorUserAlreadyLogged, correlationId, writer)
			} else {
				if t.users[login.Username()] != nil {
					t.DispatchEvent(fmt.Sprintf("User %s reconnected", login.Username()), false, 1)
				} else {
					t.DispatchEvent(fmt.Sprintf("New User %s logged in", login.Username()), false, 1)
					t.users[login.Username()] = NewUser(login.Username(), t.chEvents)
				}
				user = t.users[login.Username()]
				lastSendError = t.sendBackResponse(chat.ResponseCodeOk, correlationId, writer)
				user.UpdateWriter(writer)
			}

		case chat.CommandMessageKey:
			message := &chat.CommandMessage{}
			err := message.Read(readerFull)
			if err != nil {
				t.DispatchEvent(fmt.Sprintf("Error reading message: %v", err), true, 3)
				break
			}
			correlationId = message.CorrelationId()
			if t.users[message.To] != nil {
				t.DispatchEvent(fmt.Sprintf("Message from %s to %s: %s", message.From, message.To, message.Message), false, 2)
				lastSendError = t.sendBackResponse(chat.ResponseCodeOk, correlationId, writer)
				toUser := t.users[message.To]
				toUser.AddMessage(message.From, message.To, message.Message, message.Time)
			} else {
				t.DispatchEvent(fmt.Sprintf("User %s not found", message.To), true, 3)
				lastSendError = t.sendBackResponse(chat.ResponseCodeErrorUserNotFound, correlationId, writer)
			}
		}

		if lastSendError != nil {
			t.DispatchEvent(fmt.Sprintf("Error sending response: %v", lastSendError), true, 3)
			break
		}

		if user != nil {
			t.DispatchEvent(fmt.Sprintf("Response sent to user %s correlationId %d", user.Username, correlationId), false, 1)
		}

	}
	if user != nil {
		user.SetOnline(false)
		t.DispatchEvent(fmt.Sprintf("User %s logged out", user.Username), false, 2)
	}

}

func (t *TcpServer) sendBackResponse(code uint16, correlationId uint32, writer *bufio.Writer) error {
	genericResponse := chat.NewGenericResponse(code)
	genericResponse.SetCorrelationId(correlationId)
	return chat.WriteCommandWithHeader(genericResponse, writer)
}

func (t *TcpServer) Users() map[string]*User {
	return t.users
}
