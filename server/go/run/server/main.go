package main

import (
	"fmt"
	"gsantomaggio/chat/server/tcp_server"
	"os"
)

func main() {

	events := make(chan *tcp_server.Event)

	go func() {
		for event := range events {
			fmt.Printf("%s: %s\n", event.Time().Format("2006-01-02 15:04:05"), event.Message())
			if event.IsAnError() {
				fmt.Fprintf(os.Stderr, "error: %s\n", event.Message())
			}
		}
	}()

	tcpServer := tcp_server.NewTcpServer("localhost", 5555, events)
	err := tcpServer.StartInAThread()

	if err != nil {
		fmt.Fprintf(os.Stderr, "error starting TCP server: %v\n", err)
		return
	}

	fmt.Printf("press enter to stop the server\n")
	fmt.Scanln()
	tcpServer.Stop()

}
