package main

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"gsantomaggio/chat/server/chat"
	"gsantomaggio/chat/server/tcp_client"
	"os"
	"sync"
	"time"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <server_address>\n", os.Args[0])
		return
	}
	in := bufio.NewReader(os.Stdin)
	chMessages := make(chan *chat.CommandMessage)

	go func() {
		totalReceived := 0
		for {
			msg := <-chMessages
			totalReceived++
			color.Green("****** New message received ******\n")
			color.Green("%s -From : %s Text: %s - total: %d \n", chat.ConvertUint64ToTimeFormatted(msg.Time),
				msg.From, msg.Message, totalReceived)
			color.Green("****** End message received ******\n")
		}
	}()

	client := tcp_client.NewChatClient(chMessages)
	err := client.Connect(args[1])
	if err != nil {
		fmt.Printf("Error connecting to server: %v\n", err)
		return
	}

	fmt.Printf("Enter your user name:\n")
	username, _ := in.ReadString('\n')
	username = username[:len(username)-1]

	res, err := client.Login(username)
	if err != nil {
		return
	}

	if res.ResponseCode() != chat.ResponseCodeOk {
		fmt.Printf("Login error: %s\n", chat.FormResponseCodeToString(res.ResponseCode()))
		return
	}

	fmt.Printf("Login %s\n", chat.FormResponseCodeToString(res.ResponseCode()))

	for {
		time.Sleep(500 * time.Millisecond)
		fmt.Printf("****Menu****\n")
		fmt.Printf("1. Send a message\n")
		fmt.Printf("2. Test correlation id\n")
		fmt.Printf("3. Exit\n")
		fmt.Printf("************\n")

		option, _ := in.ReadString('\n')
		option = option[:len(option)-1]
		if option == "3" {
			break
		}

		if option == "1" {
			fmt.Printf("Write a me message to:\n")
			userTo, _ := in.ReadString('\n')
			userTo = userTo[:len(userTo)-1]
			fmt.Printf("Message text:\n")
			message, _ := in.ReadString('\n')
			message = message[:len(message)-1]
			res, err = client.SendMessage(message, userTo)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error sending message: %v\n", err)
				return
			}
			if res.ResponseCode() != chat.ResponseCodeOk {
				fmt.Fprintf(os.Stderr, "error sending message: %s\n", chat.FormResponseCodeToString(res.ResponseCode()))
			} else {

				fmt.Printf("Message sent. Response code: %s\n", chat.FormResponseCodeToString(res.ResponseCode()))
			}
		}

		if option == "2" {
			waitGroup := sync.WaitGroup{}
			for i := 0; i < 20; i++ {
				waitGroup.Add(1)
				go func(idx int) {
					fmt.Printf("Sending correlation id test\n")
					res, err = client.CorrelationIdTest()
					if err != nil {
						fmt.Fprintf(os.Stderr, "error sending message: %v\n", err)
						return
					}
					if res.ResponseCode() != chat.ResponseCodeOk {
						fmt.Fprintf(os.Stderr, "error sending message: %s\n", chat.FormResponseCodeToString(res.ResponseCode()))
					} else {
						fmt.Printf("CorrelationIdTest. Response code: %s for id %d \n", chat.FormResponseCodeToString(res.ResponseCode()), res.CorrelationId())
					}
					waitGroup.Done()
				}(i)
			}
			waitGroup.Wait()

		}
	}
}
