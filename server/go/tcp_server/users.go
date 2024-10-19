package tcp_server

import (
	"net"
	"time"
)

type UserMessage struct {
	From    string
	To      string
	Message string
	Sent    time.Time
}

type User struct {
	Username   string
	LastLogin  time.Time
	Connection net.Conn
	isOnline   bool
	Messages   []*UserMessage
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
