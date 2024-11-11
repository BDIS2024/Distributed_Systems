package main

import (
	"Lecture_7_Exercise/proto"
	"context"
	"io"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:5050", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln(err)
	}

	client := proto.NewServiceClient(conn)

	stream, err := client.DistributedMutexService(context.Background())
	if err != nil {
		log.Fatal(err.Error())
	}
	in, err := stream.Recv()
	if err == io.EOF {
		waitc <- true
		return
	}
	if err != nil {
		log.Fatal("Error receiving message:", err)
		waitc <- true
		return
	}
}

func stuff() {
	// On initialisation do
	//	 state := RELEASED;
	//
	// End on

	// On enter do
	//
	//	state := WANTED;
	//	“multicast ‘req(T,p)’”, where T := LAMPORT time of ‘req’ at p
	//	wait for N-1 replies
	//	state := HELD;
	//
	// End on

	// On receive ‘req (Ti,pi)’do
	//
	//	if(state == HELD ||
	//	(state == WANTED &&
	//	(T,pme) < (Ti,pi)))
	//	then queue req
	//	else reply to req
	//
	// End on

	// On exit do
	//
	//	state := RELEASED
	//	reply to all in queue
	//
	// End on
}
