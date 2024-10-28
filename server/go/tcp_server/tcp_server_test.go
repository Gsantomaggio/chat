package tcp_server

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gsantomaggio/chat/server/chat"
	"gsantomaggio/chat/server/tcp_client"
	"time"
)

const port = int(6666)
const address = "localhost:6666"

var _ = Describe("Tcp Server", func() {
	var tcpServer *TcpServer
	BeforeEach(func() {

		tcpServer = NewTcpServer("localhost", port, nil)
		err := tcpServer.StartInAThread()
		if err != nil {
			return
		}
		time.Sleep(200 * time.Millisecond)
	})
	AfterEach(func() {
		tcpServer.Stop()
	})
	Context("Login", func() {

		It("Login should success", func() {
			receiver := make(chan *chat.CommandMessage)
			client := tcp_client.NewChatClient(receiver)
			Expect(client.Connect(address)).To(Succeed())
			r, e := client.Login("user1")
			Expect(e).To(BeNil())
			Expect(r.ResponseCode()).To(Equal(chat.ResponseCodeOk))
			Expect(client.Close()).To(Succeed())
		})

		It("Two Logins the second should raise an error", func() {
			receiver := make(chan *chat.CommandMessage)
			client := tcp_client.NewChatClient(receiver)
			Expect(client.Connect(address)).To(Succeed())
			r, e := client.Login("user1")
			Expect(e).To(BeNil())
			Expect(r.ResponseCode()).To(Equal(chat.ResponseCodeOk))
			r, e = client.Login("user1")
			Expect(e).To(BeNil())
			Expect(r.ResponseCode()).To(Equal(chat.ResponseCodeErrorUserAlreadyLogged))
			Expect(client.Close()).To(Succeed())
			client = tcp_client.NewChatClient(receiver)
			Expect(client.Connect(address)).To(Succeed())
			r, e = client.Login("user1")
			Expect(e).To(BeNil())
			Expect(r.ResponseCode()).To(Equal(chat.ResponseCodeOk))
			Expect(client.Close()).To(Succeed())
		})
		It("Exchange Messages between two clients", func() {
			done := make(chan bool)
			receiver1 := make(chan *chat.CommandMessage)
			go func() {
				for msg := range receiver1 {
					Expect(msg).NotTo(BeNil())
					Expect(msg.From).To(Equal("user2"))
					Expect(msg.Message).To(Equal("Hello"))
					done <- true
				}
			}()
			client1 := tcp_client.NewChatClient(receiver1)
			Expect(client1.Connect(address)).To(Succeed())
			r, e := client1.Login("user1")
			Expect(e).To(BeNil())
			Expect(r.ResponseCode()).To(Equal(chat.ResponseCodeOk))
			receiver2 := make(chan *chat.CommandMessage)
			client2 := tcp_client.NewChatClient(receiver2)
			Expect(client2.Connect(address)).To(Succeed())
			r, e = client2.Login("user2")
			Expect(e).To(BeNil())
			Expect(r.ResponseCode()).To(Equal(chat.ResponseCodeOk))
			r, e = client2.SendMessage("Hello", "user1")
			Expect(e).To(BeNil())

			<-done
			close(receiver1)
			close(receiver2)
			Expect(client1.Close()).To(Succeed())
			Expect(client2.Close()).To(Succeed())
		})

	})
})
