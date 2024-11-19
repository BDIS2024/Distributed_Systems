package main

import (
	proto "Assignment_5/proto"
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"
)

type AuctionServer struct {
	proto.UnimplementedAuctionServiceServer
}

type HighestBid struct {
	mu     sync.Mutex
	value  int64
	bidder string
}

var hb HighestBid

func (s *AuctionServer) Bid(arg1 context.Context, givenBid *proto.Amount) (*proto.Ack, error) {
	fmt.Printf("Recived bid from: %s at %s\n", givenBid.Bidder, givenBid.BidTime)

	message_time, err := time.Parse(time.ANSIC, givenBid.BidTime)
	if err != nil {
		return &proto.Ack{Outcome: "Unsupported Time Format"}, nil
	}

	// Not auctioning failure state
	if !is_auctioning {
		beginAuction(message_time)
	} else if message_time.Before(end_time) {
		return &proto.Ack{Outcome: "No Longer Auctioning"}, nil
	}

	// Comparrison logic //fmt.Printf("Judgment: %v > %v = %v\n", givenBid.Bid, hb.value, givenBid.Bid > hb.value)
	hb.mu.Lock()
	if givenBid.Bid > hb.value {

		hb.value = givenBid.Bid
		hb.bidder = givenBid.Bidder

		hb.mu.Unlock()
		return &proto.Ack{Outcome: "Success"}, nil
	} else {
		hb.mu.Unlock()
		return &proto.Ack{Outcome: "Denied - Value lower than highest bidder"}, nil
	}

}

func (s *AuctionServer) Result(context.Context, *proto.Empty) (*proto.Outcome, error) {
	if !is_auctioning {
		return &proto.Outcome{HighestBid: hb.value, HighestBidder: hb.bidder, Status: "Auction Not Started"}, nil
	}

	if time.Now().Before(end_time) {
		return &proto.Outcome{HighestBid: hb.value, HighestBidder: hb.bidder, Status: "Auction Ended"}, nil
	} else {
		return &proto.Outcome{HighestBid: hb.value, HighestBidder: hb.bidder, Status: "Auction Ongoing"}, nil
	}

	//fmt.Println("Result")
}

func main() {
	hb = HighestBid{value: 0, bidder: "No Bidder"}
	is_auctioning = false

	// Setup Server

	grpcServer := grpc.NewServer()
	proto.RegisterAuctionServiceServer(grpcServer, &AuctionServer{})

	port := getPort()

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("Server started listening on port %s\n", port)
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal(err)
	}
}

var is_auctioning bool
var end_time time.Time

func beginAuction(timetostart time.Time) {
	fmt.Printf("!!!Auction started at %v!!!\n", time.Now())
	is_auctioning = true
	timeAdd, err := time.ParseDuration("100s")
	if err != nil {
		fmt.Println("Oopsie i hardcoded code?")
	}

	end_time = timetostart.Add(timeAdd)
	go alert(timeAdd)
}

func alert(sleepTime time.Duration) {

	time.Sleep(sleepTime)
	fmt.Printf("!!!Auction ended at %v !!!\nHighest Bidder: %v\n", time.Now(), hb)
}

func getPort() string {
	var port string
	var err error

	// Port
	if len(os.Args) > 1 {

		port = os.Args[1]

	} else {
		fmt.Println("Enter port number:")
		reader := bufio.NewReader(os.Stdin)
		port, err = reader.ReadString('\n')
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	port = strings.TrimSpace(port)
	port = ":" + port
	return port
}

// TODO: LEADER ELECTION & REPLICATION
