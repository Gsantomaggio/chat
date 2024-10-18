package main

import (
	"fmt"
	"gsantomaggio/chat/server/pkg"
	"os"
)

func main() {

	tcpServer := pkg.NewTcpServer("localhost", 5555)
	err := tcpServer.StartInAThread()

	if err != nil {
		fmt.Fprintf(os.Stderr, "error starting TCP server: %v\n", err)
		return
	}

	fmt.Printf("press enter to stop the server\n")
	fmt.Scanln()
	tcpServer.Stop()

}
