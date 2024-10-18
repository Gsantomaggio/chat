package main

import (
	"bufio"
	"gsantomaggio/chat/server/pkg"
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

func (f *FakeClient) ReadResponse() (*pkg.GenericResponse, error) {
	reader := bufio.NewReader(f.tcpConn)
	g := &pkg.GenericResponse{}
	err := g.Read(reader)
	return g, err
}

func (f *FakeClient) Close() error {
	return f.tcpConn.Close()
}

func (f *FakeClient) Login(login *pkg.CommandLogin) error {
	return pkg.WriteCommandWithHeader(login, bufio.NewWriter(f.tcpConn))
}

func main() {

	client := &FakeClient{}
	err := client.Connect()
	if err != nil {
		return
	}

	login := pkg.NewCommandLogin("user", 12)
	err = client.Login(login)
	if err != nil {
		return
	}

	resp, err := client.ReadResponse()
	if err != nil {
		return
	}

	if resp.ResponseCode() != pkg.ResponseCodeOk {
		return
	}

}
