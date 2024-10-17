package main

import (
	proto "Assignment_3/proto"
	"log"
	"math/rand"
	"net"
	"sync"

	"google.golang.org/grpc"
)

type ChittyChatServer struct {
	proto.UnimplementedChittyChatServiceServer
}

type messageObject struct {
	ClientName string
	Message    string
	ClientID   int
	MessageID  int
}

type messageHandler struct {
	Messages []messageObject
	Lock     sync.Mutex
}

var handler = messageHandler{}

func (s *ChittyChatServer) ChatService(stream proto.ChittyChatService_ChatServiceServer) error {
	clientID := rand.Intn(1e6)
	errorchan := make(chan error)

	go retrieveMessage(clientID, stream, errorchan)
	go sendMessage(clientID, stream, errorchan)

	return <-errorchan
}

func retrieveMessage(clientID int, stream proto.ChittyChatService_ChatServiceServer, errorchan chan error) {
	for {
		message, err := stream.Recv()
		if err != nil {
			log.Print(err.Error())
			errorchan <- err
		}
		messageID := rand.Intn(1e6)

		handler.Lock.Lock()
		handler.Messages = append(handler.Messages, messageObject{
			ClientName: message.Name,
			Message:    message.Message,
			ClientID:   clientID,
			MessageID:  messageID,
		})

		handler.Lock.Unlock()

	}
}

func sendMessage(clientID int, stream proto.ChittyChatService_ChatServiceServer, errorchan chan error) {
	for {
		for {
			handler.Lock.Lock()

			if len(handler.Messages) == 0 {
				handler.Lock.Unlock()
				break
			}

			senderID := handler.Messages[0].ClientID
			senderName := handler.Messages[0].ClientName
			senderMessage := handler.Messages[0].Message

			handler.Lock.Unlock()

			if senderID != clientID {
				err := stream.Send(&proto.ServerMessage{
					Name:      senderName,
					Message:   senderMessage,
					Timestamp: "1",
				})
				if err != nil {
					log.Println(err)
					errorchan <- err
				}

				handler.Lock.Lock()

				if len(handler.Messages) > 1 {
					handler.Messages = handler.Messages[1:]
				} else {
					handler.Messages = []messageObject{}
				}

				handler.Lock.Unlock()
			}
		}
	}
}

func main() {
	grpcServer := grpc.NewServer()
	proto.RegisterChittyChatServiceServer(grpcServer, &ChittyChatServer{})

	listener, err := net.Listen("tcp", ":5050")
	if err != nil {
		log.Fatal(err.Error())
	}

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal(err.Error())
	}
}
