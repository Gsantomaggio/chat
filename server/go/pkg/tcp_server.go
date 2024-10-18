package pkg

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

type TcpServerer interface {
}

type User struct {
	Username   string
	LastLogin  time.Time
	Connection net.Conn
	isOnline   bool
}

func NewUser(username string, Connection net.Conn) *User {
	return &User{
		Username:   username,
		LastLogin:  time.Now(),
		Connection: Connection,
		isOnline:   true,
	}
}

func (u *User) SetOnline(online bool) {
	u.isOnline = online
}

func (u *User) IsOnLine() bool {
	return u.isOnline
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
		header := &ChatHeader{}
		err := header.Read(reader)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading header: %v\n", err)
			break
		}

		fmt.Printf("Received header: %+v\n", header)

		switch header.Key() {
		case CommandLoginKey:
			login := &CommandLogin{}
			err := login.Read(reader)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error reading login: %v\n", err)
				break
			}
			fmt.Printf("Received login: %+v\n", login)
			genericResponse := NewGenericResponse(login.CorrelationId(), ResponseCodeOk)
			user = NewUser(login.Username(), conn)
			t.users[login.Username()] = user
			_, err = genericResponse.Write(writer)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error writing response: %v\n", err)
				break
			}
			err = writer.Flush()
			if err != nil {
				fmt.Fprintf(os.Stderr, "error flushing writer: %v\n", err)
				break
			}

		case CommandMessageKey:
			// handle command 2
		}

	}
	if user != nil {
		user.SetOnline(false)
	}

}

func (t *TcpServer) Users() map[string]*User {
	return t.users
}
