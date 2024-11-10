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

func main() {
	f, err := os.OpenFile("../logs", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	grpcServer := grpc.NewServer()
	proto.RegisterDmutexServiceServer(grpcServer, &DmutexServer{})

	fmt.Println("Enter port number:")
	reader := bufio.NewReader(os.Stdin)
	port, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err.Error())
	}
	port = strings.TrimSpace(port)
	port = ":" + port

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Server started on :%s", port)
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal(err)
	}
}

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

/*var handler = MessageHandler{
	Clients: make(map[string]proto.DmutexService_DmutexServer),
}*/

func (s *DmutexServer) ChatService(stream proto.DmutexService_DmutexServer) error {
	errorChan := make(chan error)

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
		fmt.Printf("Recived message: %v\n", message)
		if message.Name == "Connect" {

			// Connect
			clientNodePair = stream
			fmt.Printf("Formed a pair with:%v\n", clientNodePair)
			sendStoredMessages()

		} else {
			// redirect message to main server
			if clientNodePair == nil {
				fmt.Println("Recived message without pair - Storing message for later...")
			} else {
				sendMessageToPair(message)
			}
		}
	}
}

func sendStoredMessages() {
	if len(messageStorage) > 0 {
		for i := 0; i < len(messageStorage); i++ {
			sendMessageToPair(&messageStorage[i])
		}
		messageStorage = []proto.Message{}
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

/*

	var inCriticalSection = false
var replies = 0
		if message.Name == "Request" {


			if inCriticalSection {
				// Add client to "reply list"
			}
		} else if message.Name == "Reply" {
			// Send message to
		}*/
