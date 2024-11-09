package main

import (
	proto "Assignment_4-real-real/proto"
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
	//logs
	f, err := os.OpenFile("../logs", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	//connection
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Enter port number:")
	port, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err.Error())
	}
	port = strings.TrimSpace(port)
	host := "localhost:" + port

	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err.Error())
	}
	defer conn.Close()

	//stream
	client := proto.NewChittyChatServiceClient(conn)

	stream, err := client.ChatService(context.Background())
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Enter your name:")
	name, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err.Error())
	}
	name = strings.TrimSpace(name)

	msg := proto.ClientMessage{
		Name:      name,
		Message:   "has joined the chat.",
		Timestamp: counter,
	}

	err = stream.Send(&msg)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("Client sent request: Name: %s, Message: %s, Timestamp: (%d)\n", msg.Name, msg.Message, counter)
	waitc := make(chan bool)
	//donec := make(chan bool)

	//go retrieveMessage(waitc, donec, stream)
	//go sendMessage(donec, stream, msg.Name)

	<-waitc

}

func retrieveMessage(waitc chan bool, donec chan bool, stream proto.ChittyChatService_ChatServiceClient) {
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
				log.Fatal("Error receiving message:", err)
				waitc <- true
				return
			}
			counter = max(counter, in.Timestamp) + 1

			log.Printf("Client recieved response: Name: %s, Message: %s, Timestamp: (%d) at: %d\n", in.Name, in.Message, in.Timestamp, counter)
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
			err = stream.Send(&proto.ClientMessage{
				Name:      username,
				Message:   "has left the chat.",
				Timestamp: counter,
			})
			if err != nil {
				log.Fatal(err.Error())
			}
			err = stream.CloseSend()
			if err != nil {
				log.Fatal("Error closing stream:", err)
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
			log.Fatal("Error sending message:", err)
		}
		log.Printf("Client sent request: Name: %s, Message: %s, Timestamp: (%d)\n", msg.Name, msg.Message, counter)
	}
}

func max(counter int32, comparecounter int32) int32 {
	if counter < comparecounter {
		return comparecounter
	}
	return counter
}
