package pkg

import (
	"bufio"
	. "github.com/onsi/ginkgo/v2"
	"net"
)

type FakeClient struct {
	tcpConn *net.TCPConn
}

func (f *FakeClient) Connect() error {
	servAddr := "localhost:5555"
	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
	if err != nil {
		return err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return err
	}
	f.tcpConn = conn
	return nil
}

func (f *FakeClient) ReadResponse() (*GenericResponse, error) {
	reader := bufio.NewReader(f.tcpConn)
	g := &GenericResponse{}
	err := g.Read(reader)
	return g, err
}

func (f *FakeClient) Close() error {
	return f.tcpConn.Close()
}

func (f *FakeClient) Login(login *CommandLogin) error {
	return WriteCommandWithHeader(login, bufio.NewWriter(f.tcpConn))
}

var _ = Describe("Tcp Server", func() {
	Context("Login", func() {
		//var tcpServer *TcpServer
		BeforeEach(func() {

			//tcpServer = NewTcpServer("localhost", 5555)
			//err := tcpServer.StartInAThread()
			//if err != nil {
			//	return
			//}
		})
		AfterEach(func() {
			//tcpServer.Stop()
		})
		It("Login should success", func() {
			//client := &FakeClient{}
			//err := client.Connect()
			//if err != nil {
			//	return
			//}
			//login := NewCommandLogin("user", 1)
			//err = client.Login(login)
			//if err != nil {
			//	return
			//}
			//_, err = client.ReadResponse()
			//if err != nil {
			//	return
			//}

		})
	})
})
