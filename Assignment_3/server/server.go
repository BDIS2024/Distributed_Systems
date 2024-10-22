package main

import (
	proto "Assignment_3/proto"
	"io"
	"log"
	"net"
	"os"
	"sync"

	"google.golang.org/grpc"
)

type ChittyChatServer struct {
	proto.UnimplementedChittyChatServiceServer
}

type MessageObject struct {
	ClientName string
	Message    string
	Timestamp  int32
}

type MessageHandler struct {
	Clients map[string]proto.ChittyChatService_ChatServiceServer
	Lock    sync.Mutex
}

var handler = MessageHandler{
	Clients: make(map[string]proto.ChittyChatService_ChatServiceServer),
}

var counter int32 = 0

func (s *ChittyChatServer) ChatService(stream proto.ChittyChatService_ChatServiceServer) error {
	errorChan := make(chan error)

	go retrieveMessagesFromClient(stream, errorChan)

	return <-errorChan
}

func retrieveMessagesFromClient(stream proto.ChittyChatService_ChatServiceServer, errorChan chan error) {
	clientName := ""

	for {
		message, err := stream.Recv()
		if err == io.EOF {
			removeClient(clientName)
			errorChan <- err
			return
		}
		if err != nil {
			log.Printf("Error receiving message: %v", err)
			removeClient(clientName)
			errorChan <- err
			return
		}

		if clientName == "" {
			clientName = message.Name
			addClient(clientName, stream)
		}

		if len(message.Message) > 128 {
			sendErrorToCLient(clientName, "Message has to be under 128 characters.")
			continue
		}

		counter = max(counter, message.Timestamp) + 1
		log.Printf("Server recieved request: Name: %s, Message: %s, Timestamp: (%d) at %d\n", message.Name, message.Message, message.Timestamp, counter)
		broadcastMessageToClients(message)
	}
}

func broadcastMessageToClients(message *proto.ClientMessage) {
	handler.Lock.Lock()
	defer handler.Lock.Unlock()
	counter++
	for clientName, clientStream := range handler.Clients {
		err := clientStream.Send(&proto.ServerMessage{
			Name:      message.Name,
			Message:   message.Message,
			Timestamp: counter,
		})
		if err != nil {
			log.Printf("Error sending message to %s: %v", clientName, err)
			removeClient(clientName)
		}
	}
	log.Printf("Server sent response: Name: %s, Message: %s, Timestamp: (%d)\n", message.Name, message.Message, counter)
}

func sendErrorToCLient(clientName string, erro string) {
	handler.Lock.Lock()
	defer handler.Lock.Unlock()
	counter++

	err := handler.Clients[clientName].Send(&proto.ServerMessage{
		Name:      "Server",
		Message:   erro,
		Timestamp: counter,
	})
	if err != nil {
		log.Printf("Error sending message to %s: %v", clientName, err)

	}
}

func addClient(clientName string, client proto.ChittyChatService_ChatServiceServer) {
	handler.Lock.Lock()
	defer handler.Lock.Unlock()

	handler.Clients[clientName] = client
	counter++
	log.Printf("Server added client: %s, (%d)\n", clientName, counter)
}

func removeClient(clientName string) {
	handler.Lock.Lock()
	defer handler.Lock.Unlock()

	delete(handler.Clients, clientName)
	counter++
	log.Printf("Server removed client: %s, (%d)\n", clientName, counter)
}

func main() {
	f, err := os.OpenFile("../logs", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	grpcServer := grpc.NewServer()
	proto.RegisterChittyChatServiceServer(grpcServer, &ChittyChatServer{})

	listener, err := net.Listen("tcp", ":5050")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Server started on :5050")
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal(err)
	}
}

func max(counter int32, comparecounter int32) int32 {
	if counter < comparecounter {
		return comparecounter
	}
	return counter
}
