package main

import (
	"fmt"
	"github.com/fatih/color"
	"gsantomaggio/chat/server/tcp_server"
	"os"
)

func printColoredMessage(event *tcp_server.Event) {
	switch event.Level() {
	case 1:
		color.Cyan("%s: %s\n", event.Time().Format("2006-01-02 15:04:05"), event.Message())
	case 2:
		color.Green("%s: %s\n", event.Time().Format("2006-01-02 15:04:05"), event.Message())
	case 3:
		color.Red("%s: %s\n", event.Time().Format("2006-01-02 15:04:05"), event.Message())
	case 4:
		color.Yellow("%s: %s\n", event.Time().Format("2006-01-02 15:04:05"), event.Message())
	default:
		color.White("%s: %s\n", event.Time().Format("2006-01-02 15:04:05"), event.Message())
	}

}

func main() {

	events := make(chan *tcp_server.Event)

	go func() {
		for event := range events {
			printColoredMessage(event)
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
