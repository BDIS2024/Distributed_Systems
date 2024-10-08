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
	messages []string
}

func (s *server) GetMesssages(ctx context.Context, in *proto.Empty) (*proto.Messages, error) {
	return &proto.Messages{Messages: s.messages}, nil
}

func (s *server) SendMessage(ctx context.Context, in *proto.Message) (*proto.Empty, error) {
	server := &server{}
	s.messages = append(server.messages, in.Message)
	return &proto.Empty{}, nil
}

func main() {
	server := &server{messages: []string{}}
	server.messages = append(server.messages, "hello")
	server.messages = append(server.messages, "hello2")

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
