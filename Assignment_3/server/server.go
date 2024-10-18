package main

import (
	proto "Assignment_3/proto"
	"io"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
)

type ChittyChatServer struct {
	proto.UnimplementedChittyChatServiceServer
}

type MessageObject struct {
	ClientName string
	Message    string
}

type MessageHandler struct {
	Clients map[string]proto.ChittyChatService_ChatServiceServer
	Lock    sync.Mutex
}

var handler = MessageHandler{
	Clients: make(map[string]proto.ChittyChatService_ChatServiceServer),
}

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

		broadcastMessage(MessageObject{
			ClientName: message.Name,
			Message:    message.Message,
		})
	}
}

func broadcastMessage(message MessageObject) {
	handler.Lock.Lock()
	defer handler.Lock.Unlock()

	for clientName, clientStream := range handler.Clients {

		if clientName != message.ClientName {
			err := clientStream.Send(&proto.ServerMessage{
				Name:      message.ClientName,
				Message:   message.Message,
				Timestamp: "1",
			})
			if err != nil {
				log.Printf("Error sending message to %s: %v", clientName, err)
				removeClient(clientName)
			}
		}
	}
}

func addClient(clientName string, client proto.ChittyChatService_ChatServiceServer) {
	handler.Lock.Lock()
	defer handler.Lock.Unlock()

	handler.Clients[clientName] = client
}

func removeClient(clientName string) {
	handler.Lock.Lock()
	defer handler.Lock.Unlock()

	delete(handler.Clients, clientName)
}

func main() {
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
