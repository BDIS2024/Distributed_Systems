package main

import (
	proto "Assignment_3/proto"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
)

type ChittyChatServer struct {
	proto.UnimplementedChittyChatServiceServer
}

func (s *ChittyChatServer) ChatService(stream proto.ChittyChatService_ChatServiceServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatal(err.Error())
		}

		stream.Send(in)
	}
}

func main() {
	server := &ChittyChatServer{}
	server.start_server()
}

func (s *ChittyChatServer) start_server() {
	grpcServer := grpc.NewServer()
	proto.RegisterChittyChatServiceServer(grpcServer, s)

	listener, err := net.Listen("tcp", ":5050")
	if err != nil {
		log.Fatalf("Did not work")
	}

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("Did not work")
	}

}
