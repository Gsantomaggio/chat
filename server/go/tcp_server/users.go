package tcp_server

import (
	"bufio"
	"fmt"
	"gsantomaggio/chat/server/chat"
	"net"
	"sync"
	"time"
)

type UserMessage struct {
	From    string
	To      string
	Message string
	Sent    uint64
}

type User struct {
	Username   string
	LastLogin  time.Time
	Connection net.Conn
	isOnline   bool
	Messages   []*UserMessage
	chNotify   chan struct{}
	mutex      sync.Mutex
	chEvents   chan *Event
	writer     *bufio.Writer
}

func NewUser(username string, chEvents chan *Event) *User {
	u := &User{
		Username:  username,
		LastLogin: time.Now(),
		isOnline:  true,
		Messages:  make([]*UserMessage, 0),
		chNotify:  make(chan struct{}),
		mutex:     sync.Mutex{},
		chEvents:  chEvents,
	}
	u.sendMessageInAThread()
	return u
}

func (u *User) SetOnline(online bool) {
	u.isOnline = online
	if online {
		u.chNotify <- struct{}{}
	}
}

func (u *User) UpdateWriter(writer *bufio.Writer) {
	u.writer = writer
	u.SetOnline(true)
}

func (u *User) IsOnLine() bool {
	return u.isOnline
}

func (u *User) Close() {
	u.SetOnline(false)
}

func (u *User) DispatchEvent(message string, isAnError bool) {
	if u.chEvents != nil {
		u.chEvents <- NewEvent(message, isAnError)
	}
}

func (u *User) AddMessage(from, to, message string, sent uint64) {
	u.mutex.Lock()
	u.Messages = append(u.Messages, &UserMessage{
		From:    from,
		To:      to,
		Message: message,
		Sent:    sent,
	})
	u.mutex.Unlock()
	if u.isOnline {
		u.chNotify <- struct{}{}
	} else {
		u.DispatchEvent(fmt.Sprintf("User %s is offline and received a message from %s", u.Username, from), false)
	}
}

func (u *User) sendMessageInAThread() {
	go func() {
		for _ = range u.chNotify {
			u.mutex.Lock()
			for _, message := range u.Messages {
				//time.Sleep(100 * time.Millisecond)
				if message.To != u.Username {
					u.DispatchEvent(fmt.Sprintf("Message from %s to %s not sent", message.From, u.Username), false)
					continue
				}
				err := chat.WriteCommandWithHeader(
					chat.NewCommandMessageWithCorrelationId(
						message.Message,
						message.From, u.Username,
						0, message.Sent),
					u.writer)
				u.DispatchEvent(fmt.Sprintf("Sent message from %s to %s message: %s", message.From, u.Username, message.Message), false)
				if err != nil {
					u.DispatchEvent(fmt.Sprintf("Error sending message to %s: %v", u.Username, err), true)
					break
				}
			}
			u.Messages = make([]*UserMessage, 0)
			u.mutex.Unlock()
		}
	}()
}
