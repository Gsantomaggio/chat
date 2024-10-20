package main

import (
	"bufio"
	"fmt"
	"gsantomaggio/chat/server/chat"
	"gsantomaggio/chat/server/tcp_client"
	"os"
)

func main() {

	chMessages := make(chan *chat.CommandMessage)

	go func() {
		totalReceived := 0
		for {
			msg := <-chMessages
			totalReceived++
			fmt.Printf("%s - From : %s Text: %s - total: %d \n", chat.ConvertUint64ToTimeFormatted(msg.Time),
				msg.From, msg.Message, totalReceived)
		}
	}()

	client := tcp_client.NewChatClient(chMessages)
	err := client.Connect("localhost:5555")
	if err != nil {
		fmt.Printf("Error connecting to server: %v\n", err)
		return
	}

	in := bufio.NewReader(os.Stdin)
	println("Insert your Username:")
	username, err := in.ReadString('\n')
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
		fmt.Printf("Destination:\n")
		userTo, _ := in.ReadString('\n')
		userTo = userTo[:len(userTo)-1]
		fmt.Printf("Message:\n")
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

}
