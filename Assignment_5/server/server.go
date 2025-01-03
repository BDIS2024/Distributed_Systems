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
var port string

func (s *AuctionServer) Bid(arg1 context.Context, givenBid *proto.Amount) (*proto.Ack, error) {
	log.Printf("Received bid from: %s at %s\n", givenBid.Bidder, givenBid.BidTime)
	fmt.Printf("Received bid from: %s at %s\n", givenBid.Bidder, givenBid.BidTime)

	message_time, err := time.Parse(time.RFC850, givenBid.BidTime)
	if err != nil {
		return &proto.Ack{Outcome: "Unsupported Time Format"}, nil
	}

	// Not auctioning failure state
	if !is_auctioning {
		beginAuction(message_time)
	} else if message_time.After(end_time) {
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

	tmp := time.Now()
	//fmt.Printf("Judgment: %v > %v = %v\n", tmp, end_time, tmp.After(end_time))
	if tmp.After(end_time) {
		return &proto.Outcome{HighestBid: hb.value, HighestBidder: hb.bidder, Status: "Auction Ended"}, nil
	} else {
		return &proto.Outcome{HighestBid: hb.value, HighestBidder: hb.bidder, Status: "Auction Ongoing"}, nil
	}

	//fmt.Println("Result")
}

func main() {
	//logs
	f, err := os.OpenFile("../logs", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

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

	log.Printf("Server started listening on port %s\n", port)
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
		fmt.Println("Oopsie i hardcoded code? // what ??")
	}

	end_time = timetostart.Add(timeAdd)

	log.Printf("Port: %s... Auction supposed to end at %v\n", port, end_time)
	fmt.Printf("Auction supposed to end at %v\n", end_time)
	go alert(timeAdd)
}

func alert(sleepTime time.Duration) {

	time.Sleep(sleepTime)
	log.Printf("Port: %s... !!!Auction ended at %v !!! Highest Bidder: %v, with a bid of %v\n", port, time.Now(), hb.bidder, hb.value)
	fmt.Printf("!!!Auction ended at %v !!!\nHighest Bidder: %v, with a bid of %v\n", time.Now(), hb.bidder, hb.value)
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
