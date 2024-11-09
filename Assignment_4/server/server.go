package main

import (
	proto "Assignment_4/proto"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
)

type dmutexserver struct {
	proto.UnimplementedMutualExclusionServer
}

func (s *dmutexserver) RequestAccess(stream proto.MutualExclusion_RequestAccessServer) error {
	errorChan := make(chan error)
	go receiveAccessRequest(stream, errorChan)

	return <-errorChan
}

func (s *dmutexserver) Release(stream proto.MutualExclusion_ReleaseServer) error {
	errorChan := make(chan error)
	go receiveReleaseRequest(stream, errorChan)

	return <-errorChan
}

func receiveAccessRequest(stream proto.MutualExclusion_RequestAccessServer, errchan chan error) {
	for {
		accessRequest, err := stream.Recv()
		if err == io.EOF {
			errchan <- err
			return
		}
		if err != nil {
			log.Printf("Error receiving message: %v\n", err)
			errchan <- err
			return
		}
		log.Printf("Received access request from: %s, at timestamp: %d\n", accessRequest.NodeId, accessRequest.Timestamp)
	}
}

func receiveReleaseRequest(stream proto.MutualExclusion_ReleaseServer, errchan chan error) {
	for {
		releaseRequest, err := stream.Recv()
		if err == io.EOF {
			errchan <- err
			return
		}
		if err != nil {
			log.Printf("Error receiving message: %v\n", err)
			errchan <- err
			return
		}
		log.Printf("Received release request from: %s\n", releaseRequest.NodeId)
	}
}

func main() {

	grpcServer := grpc.NewServer()
	proto.RegisterMutualExclusionServer(grpcServer, &dmutexserver{})

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
