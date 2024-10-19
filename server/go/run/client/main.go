package main

import (
	"fmt"
	"gsantomaggio/chat/server/chat"
	"gsantomaggio/chat/server/tcp_client"
)

func main() {

	chMessages := make(chan *chat.CommandMessage)

	go func() {
		for {
			msg := <-chMessages
			fmt.Printf("Received message: %+v\n", msg)
		}
	}()

	client := tcp_client.NewChatClient(chMessages)
	err := client.Connect("localhost:5555")
	if err != nil {
		fmt.Printf("Error connecting to server: %v\n", err)
		return
	}

	println("your username")
	var username string
	fmt.Scanf("%s", &username)

	res, err := client.Login(username)
	if err != nil {
		return
	}
	fromCodeToString := "Success"
	switch res.ResponseCode() {
	case chat.ResponseCodeErrorUserAlreadyLogged:
		fromCodeToString = "ErrorUserAlreadyLogged"
	case chat.ResponseCodeErrorUserNotFound:
		fromCodeToString = "ErrorUserNotFound"
	}

	if res.ResponseCode() != chat.ResponseCodeOk {
		fmt.Printf("Login error: %s\n", fromCodeToString)
		return
	}

	fmt.Printf("Login %s\n", fromCodeToString)

	for {
		fmt.Printf("write the user to send a message\n")
		var userTo string
		fmt.Scanf("%s", &userTo)
		fmt.Printf("write the message\n")
		var message string
		fmt.Scanf("%s", &message)

		res, err = client.SendMessage(message, userTo)
		if err != nil {
			return
		}
		if res.ResponseCode() != chat.ResponseCodeOk {
			return
		}

		fmt.Printf("Message sent\n")
	}

}
