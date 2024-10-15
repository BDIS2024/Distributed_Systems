package main

import (
	proto "Assignment_3/proto"
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
)

type ChittyChatServer struct {
	proto.UnimplementedChittyChatServiceServer
	messages []*proto.Message
}

func (s *ChittyChatServer) GetMessages(ctx context.Context, in *proto.Empty) (*proto.Messages, error) {
	return &proto.Messages{Messages: s.messages}, nil
}

func (s *ChittyChatServer) SendMessage(ctx context.Context, in *proto.Message) (*proto.Message, error) {
	s.messages = append(s.messages, in)
	return in, nil
}

func main() {
	server := &ChittyChatServer{messages: []*proto.Message{}}
	server.start_server()
}

func (s *ChittyChatServer) start_server() {
	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", ":5050")
	if err != nil {
		log.Fatalf("Did not work")
	}

	proto.RegisterChittyChatServiceServer(grpcServer, s)

	err = grpcServer.Serve(listener)

	if err != nil {
		log.Fatalf("Did not work")
	}

}
