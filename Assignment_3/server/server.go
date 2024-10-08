package main

import (
	proto "Assignment_3/proto"
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	proto.UnimplementedChittyChatServiceServer
	messages []*proto.Message
}

func (s *server) GetMesssages(ctx context.Context, in *proto.Empty) (*proto.Messages, error) {
	return &proto.Messages{Messages: s.messages}, nil
}

func (s *server) SendMessages(ctx context.Context, in *proto.Message) (*proto.Empty, error) {
	return &proto.Empty{}, nil
}

func main() {
	server := &server{messages: []*proto.Message{}}
	server.messages = append(server.messages, &proto.Message{Message: "hello", Timestamp: 1})

	server.start_server()
}

func (s *server) start_server() {
	grpcServer := grpc.NewServer()
	proto.RegisterChittyChatServiceServer(grpcServer, s)

	listener, err := net.Listen("tcp", ":5050")
	if err != nil {
		log.Fatalf("connection error")
	}
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("serve error")
	}

}
