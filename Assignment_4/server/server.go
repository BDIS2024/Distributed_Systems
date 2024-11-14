package main

import (
	proto "Assignment_4/proto"
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"

	"google.golang.org/grpc"
)

var port string

func main() {
	f, err := os.OpenFile("../logs", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	grpcServer := grpc.NewServer()
	proto.RegisterDmutexServiceServer(grpcServer, &DmutexServer{})

	port = getPort()

	fmt.Printf("Setting up listener.\n")
	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}

	fmt.Printf("Serving port\n")
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Server started on %s\n", port)
	log.Printf("Server started on %s\n", port)
}

var counter int32 = 0

type DmutexServer struct {
	proto.UnimplementedDmutexServiceServer
}

type MessageObject struct {
	ClientName string
	Message    string
	Timestamp  int32
}

type MessageHandler struct {
	Clients map[string]proto.DmutexService_DmutexServer
	Lock    sync.Mutex
}

func (s *DmutexServer) Dmutex(stream proto.DmutexService_DmutexServer) error {
	errorChan := make(chan error)
	messageStorage = []proto.Message{}

	go retrieveMessagesFromClient(stream, errorChan)

	return <-errorChan
}

var clientNodePair proto.DmutexService_DmutexServer
var messageStorage []proto.Message

func retrieveMessagesFromClient(stream proto.DmutexService_DmutexServer, errorChan chan error) {
	for {
		message, err := stream.Recv()
		if err == io.EOF {
			errorChan <- err
			return
		}
		if err != nil {
			log.Printf("Error receiving message: %v", err)
			errorChan <- err
			return
		}

		// HANDLE MESSAGE
		fmt.Printf("Server - Recived message: %v\n", message)
		log.Printf("Server: %s - Recived message: %v\n", port, message)
		var recievedTimestamp = message.Timestamp
		counter = max(counter, recievedTimestamp) + 1

		if message.Message == "Connect" {

			// Connect
			clientNodePair = stream
			fmt.Printf("Server - Formed a pair with:%v at Lamport time: %v \n", clientNodePair, counter)
			log.Printf("Server %s - Formed a pair with:%v at Lamport time: %v \n", port, clientNodePair, counter)
			sendStoredMessages()

		} else {
			// redirect message to main server
			if clientNodePair == nil {
				fmt.Printf("Server - Recived message without pair at Lamport time: %v - Storing message for later...\n", counter)
				log.Printf("Server %s - Recived message without pair at Lamport time: %v - Storing message for later...\n", port, counter)

				messageStorage = append(messageStorage, copyMessage(message))

			} else {
				sendMessageToPair(message)
			}
		}
	}
}

func sendStoredMessages() {
	var latest int32 = 0
	if len(messageStorage) > 0 {
		for i := 0; i < len(messageStorage); i++ {
			sendMessageToPair(&messageStorage[i])
			if messageStorage[i].Timestamp > latest {
				latest = messageStorage[i].Timestamp
			}
		}
		counter = max(counter, latest)
		messageStorage = []proto.Message{}
	}
}

func copyMessage(arg *proto.Message) proto.Message {
	return proto.Message{
		Name:      arg.Name,
		Message:   arg.Message,
		Timestamp: arg.Timestamp,
	}
}

func sendMessageToPair(message *proto.Message) {
	err := clientNodePair.Send(message)
	if err != nil {
		log.Printf("Error sending message: %v\n", message)
	} else {
		fmt.Printf("Sucessfully sent message: %v\n", message)
	}
}

func getPort() string {
	var port string
	var err error

	// Port
	if len(os.Args) > 1 {

		fmt.Printf("test:%v\n", os.Args[1])
		port = os.Args[1]
		if err != nil {
			log.Fatal(err.Error())
		}
	} else {
		fmt.Println("Enter port number:")
		reader := bufio.NewReader(os.Stdin)
		port, err = reader.ReadString('\n')
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	port = strings.TrimSpace(port)
	port = ":" + port
	return port
}

func max(counter int32, comparecounter int32) int32 {
	if counter < comparecounter {
		return comparecounter
	}
	return counter
}
