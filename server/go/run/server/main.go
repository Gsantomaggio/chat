package main

import (
	"fmt"
	"gsantomaggio/chat/server/tcp_server"
	"os"
)

func main() {

	tcpServer := tcp_server.NewTcpServer("localhost", 5555)
	err := tcpServer.StartInAThread()

	if err != nil {
		fmt.Fprintf(os.Stderr, "error starting TCP server: %v\n", err)
		return
	}

	fmt.Printf("press enter to stop the server\n")
	fmt.Scanln()
	tcpServer.Stop()

}
