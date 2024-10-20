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

func (t *TcpServer) DispatchEvent(message string, isAnError bool) {
	if t.chEvents != nil {
		t.chEvents <- NewEvent(message, isAnError)
	}
}

func (t *TcpServer) StartInAThread() error {
	go func() {
		err := t.Start()
		if err != nil {
			t.DispatchEvent(fmt.Sprintf("Error starting server: %v", err), true)
		}
	}()
	return nil
}
func (t *TcpServer) Start() error {
	address := fmt.Sprintf("%s:%d", t.host, t.port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		t.DispatchEvent(fmt.Sprintf("Error starting server: %v", err), true)
		return fmt.Errorf("error starting TCP server: %v", err)
	}
	t.listener = listener

	t.DispatchEvent(fmt.Sprintf("Server started at %s", address), false)
	for {
		conn, err := listener.Accept()
		if err != nil {
			t.DispatchEvent(fmt.Sprintf("Error accepting connection: %v", err), true)
			break
		}
		go t.handleConnection(conn)
	}

	t.DispatchEvent("Server stopped", false)
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
			t.DispatchEvent(fmt.Sprintf("Error reading source: %v", err), true)
			return
		}

		header := &chat.ChatHeader{}
		err = header.Read(readerFull)
		if err != nil {
			if errors.Is(err, io.EOF) {
				t.DispatchEvent("Connection closed due of EOF", false)
			} else {
				t.DispatchEvent(fmt.Sprintf("Error reading header: %v", err), true)
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
				t.DispatchEvent(fmt.Sprintf("Error reading login: %v", err), true)
				break
			}
			correlationId = login.CorrelationId()
			t.DispatchEvent(fmt.Sprintf("Login request for user %s", login.Username()), false)
			if t.users[login.Username()] != nil && t.users[login.Username()].IsOnLine() {
				t.DispatchEvent(fmt.Sprintf("User %s already logged", login.Username()), false)
				lastSendError = t.sendBackResponse(chat.ResponseCodeErrorUserAlreadyLogged, correlationId, writer)
			} else {
				if t.users[login.Username()] != nil {
					t.DispatchEvent(fmt.Sprintf("User %s reconnected", login.Username()), false)
				} else {
					t.DispatchEvent(fmt.Sprintf("New User %s logged in", login.Username()), false)
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
				t.DispatchEvent(fmt.Sprintf("Error reading message: %v", err), true)
				break
			}
			correlationId = message.CorrelationId()
			if t.users[message.To] != nil {
				t.DispatchEvent(fmt.Sprintf("Message from %s to %s: %s", message.From, message.To, message.Message), false)
				lastSendError = t.sendBackResponse(chat.ResponseCodeOk, correlationId, writer)
				toUser := t.users[message.To]
				toUser.AddMessage(message.From, message.To, message.Message, message.Time)
			} else {
				t.DispatchEvent(fmt.Sprintf("User %s not found", message.To), false)
				lastSendError = t.sendBackResponse(chat.ResponseCodeErrorUserNotFound, correlationId, writer)
			}
		}

		if lastSendError != nil {
			t.DispatchEvent(fmt.Sprintf("Error sending response: %v", lastSendError), true)
			break
		}

		if user != nil {
			t.DispatchEvent(fmt.Sprintf("Response sent to user %s correlationId %d", user.Username, correlationId), false)
		}

	}
	if user != nil {
		user.SetOnline(false)
		t.DispatchEvent(fmt.Sprintf("User %s logged out", user.Username), false)
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
