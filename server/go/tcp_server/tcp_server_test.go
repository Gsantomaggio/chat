package tcp_server

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gsantomaggio/chat/server/chat"
	"gsantomaggio/chat/server/tcp_client"
)

const port = int(6666)
const address = "localhost:6666"

var _ = Describe("Tcp Server", func() {
	var tcpServer *TcpServer
	BeforeEach(func() {

		tcpServer = NewTcpServer("localhost", port)
		err := tcpServer.StartInAThread()
		if err != nil {
			return
		}
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
	})
})
