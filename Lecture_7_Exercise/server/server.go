package main

import (
	proto "Lecture_7_Exercise/proto"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	proto.UnimplementedServiceServer
}

func main() {
	grpc := grpc.NewServer()
	proto.RegisterServiceServer(grpc, &Server{})

	listener, err := net.Listen("tcp", ":5050")
	if err != nil {
		log.Fatalln(err)
	}

	err = grpc.Serve(listener)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Server started listening on port: 5050")
}
