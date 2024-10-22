package main

import (
	proto "Assignment_3/proto"
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var counter int32 = 0

func main() {
	conn, err := grpc.NewClient("localhost:5050", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err.Error())
	}
	defer conn.Close()

	client := proto.NewChittyChatServiceClient(conn)

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Enter your username")
	username, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err.Error())
	}
	username = strings.TrimSpace(username)

	for {
		fmt.Println("Type 'join' to join a chat session.")
		fmt.Println("Type 'exit' to exit the program")

		message, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err.Error())
		}
		message = strings.TrimSpace(message)

		if message == "join" {
			consoleChannel := make(chan proto.ClientMessage) // channel for console.

			stream, err := client.ChatService(context.Background())
			if err != nil {
				log.Fatal(err.Error())
			}

			msg := proto.ClientMessage{
				Name:      username,
				Message:   "has joined the chat.",
				Timestamp: counter,
			}

			err = stream.Send(&msg)
			if err != nil {
				log.Println(err.Error())
			}

			waitc := make(chan bool)
			donec := make(chan bool)

			go retrieveMessage(waitc, donec, stream, consoleChannel)
			go sendMessage(donec, stream, username)
			go consoleManager(consoleChannel)

			<-waitc
		} else if message == "exit" {
			fmt.Println("Exiting program...")
			break
		}
	}
}

func CallClear() {
	for i := 0; i < 30; i++ {
		fmt.Println()
	}
}

func consoleManager(consoleChannel chan proto.ClientMessage) {
	messages := []proto.ClientMessage{}
	for {
		in := <-consoleChannel
		messages = append(messages, in)
		CallClear()

		fmt.Println("--- Chitty-Chat ---")

		i := 0
		if len(messages) > 20 {
			i = len(messages) - 20
			fmt.Printf("<%d Previous messages>\n", i)
		}

		for ; i < len(messages); i++ {
			msg := messages[i]
			fmt.Printf("%s : %s (%d)\n", msg.Name, msg.Message, msg.Timestamp)
		}
	}
}

func retrieveMessage(waitc chan bool, donec chan bool, stream proto.ChittyChatService_ChatServiceClient, consoleChannel chan proto.ClientMessage) {
	for {
		select {
		case <-donec:
			waitc <- true
			return
		default:
			in, err := stream.Recv()
			if err == io.EOF {
				waitc <- true
				return
			}
			if err != nil {
				log.Println("Error receiving message:", err)
				waitc <- true
				return
			}
			counter = max(counter, in.Timestamp) + 1
			consoleChannel <- proto.ClientMessage{
				Name:      in.Name,
				Message:   in.Message,
				Timestamp: counter,
			}
		}
	}
}

func sendMessage(donec chan bool, stream proto.ChittyChatService_ChatServiceClient, username string) {
	reader := bufio.NewReader(os.Stdin)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err.Error())
		}
		message = strings.TrimSpace(message)

		if message == "leave" {
			fmt.Printf("%s has left the chat. (%d)\n", username, counter)

			err = stream.Send(&proto.ClientMessage{
				Name:      username,
				Message:   "has left the chat.",
				Timestamp: counter,
			})
			if err != nil {
				log.Println(err.Error())
			}
			err = stream.CloseSend()
			if err != nil {
				log.Println("Error closing stream:", err)
			}

			donec <- true
			return
		}
		counter++

		msg := proto.ClientMessage{
			Name:      username,
			Message:   message,
			Timestamp: counter,
		}
		err = stream.Send(&msg)
		if err != nil {
			log.Println("Error sending message:", err)
		}
	}
}

func max(counter int32, comparecounter int32) int32 {
	if counter < comparecounter {
		return comparecounter
	}
	return counter
}
