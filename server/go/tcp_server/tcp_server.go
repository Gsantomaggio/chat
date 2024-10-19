package tcp_server

import (
	"bufio"
	"errors"
	"fmt"
	"gsantomaggio/chat/server/chat"
	"io"
	"net"
	"os"
)

type TcpServerer interface {
}

type TcpServer struct {
	host     string
	port     int
	users    map[string]*User
	listener net.Listener
}

func NewTcpServer(host string, port int) *TcpServer {
	return &TcpServer{
		host:  host,
		port:  port,
		users: make(map[string]*User),
	}
}

func (t *TcpServer) StartInAThread() error {
	go func() {
		err := t.Start()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error starting TCP server: %v\n", err)
		}
	}()
	return nil
}
func (t *TcpServer) Start() error {
	address := fmt.Sprintf("%s:%d", t.host, t.port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("error starting TCP server: %v", err)
	}
	t.listener = listener

	fmt.Printf("Server started on %s\n", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error accepting connection: %v\n", err)
			break
		}
		go t.handleConnection(conn)
	}

	fmt.Printf("Server stopped\n")
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
		header := &chat.ChatHeader{}
		err := header.Read(reader)
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Fprintf(os.Stderr, "Closing connection due of: %v\n", err)
			}
			break
		}
		code := chat.ResponseCodeOk
		var correlationId uint32
		switch header.Key() {
		case chat.CommandLoginKey:
			login := &chat.CommandLogin{}
			err := login.Read(reader)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error reading login: %v\n", err)
				break
			}
			correlationId = login.CorrelationId()

			if t.users[login.Username()] != nil && t.users[login.Username()].IsOnLine() {
				code = chat.ResponseCodeErrorUserAlreadyLogged
			} else {
				user = NewUser(login.Username(), conn)
				t.users[login.Username()] = user
			}

		case chat.CommandMessageKey:
			message := &chat.CommandMessage{}
			err := message.Read(reader)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error reading message: %v\n", err)
				break
			}
			correlationId = message.CorrelationId()
			if t.users[message.To] != nil {
				code = chat.ResponseCodeOk
				toUser := t.users[message.To]
				err := chat.WriteCommandWithHeader(chat.NewCommandMessageWithCorrelationId(message.Message, message.From, message.To, message.CorrelationId()),

					bufio.NewWriter(toUser.Connection))
				if err != nil {
					return
				}
			} else {
				code = chat.ResponseCodeErrorUserNotFound
				fmt.Fprintf(os.Stderr, "User %s not found\n", message.To)
			}
		}

		genericResponse := chat.NewGenericResponse(chat.GenericResponseKey, code)
		genericResponse.SetCorrelationId(correlationId)
		err = chat.WriteCommandWithHeader(genericResponse, writer)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error writing response: %v\n", err)
			return
		}

	}
	if user != nil {
		user.SetOnline(false)
	}

}

func (t *TcpServer) Users() map[string]*User {
	return t.users
}
