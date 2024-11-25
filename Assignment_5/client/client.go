package main

import (
	proto "Assignment_5/proto"
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var name string
var auctionservers []string
var auctionserverconnections []proto.AuctionServiceClient
var output []*proto.Outcome

func main() {
	//logs
	f, err := os.OpenFile("../logs", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	start := time.Now()
	auctionservers = append(auctionservers, "5050", "5051", "5052")
	connectToServers()

	name = getName()

	wait := make(chan bool)
	var checkoutput bool

	go prompt(wait)

	for {
		output = nil
		checkoutput = false
		for i := 0; i < len(auctionserverconnections); i++ {
			result, err := auctionserverconnections[i].Result(context.Background(), &proto.Empty{})
			if err != nil {
				log.Printf("Client: %v... Auctionserver %s has crashed.\n", name, auctionservers[i])
				fmt.Printf("Auctionserver %s has crashed.\n", auctionservers[i])
				auctionservers = removePort(auctionservers, i)
				auctionserverconnections = removeConn(auctionserverconnections, i)
				continue
			}
			checkoutput = true
			output = append(output, result)
		}

		if len(auctionserverconnections) == 0 {
			log.Printf("Client: %v... All auction servers are down.\n", name)
			fmt.Println("All auction servers are down.")
			wait <- true
			break
		}
		if !ongoing(output) && checkoutput {
			log.Printf("Client: %v...  Auction has ended. The highest bidder was %s with a bid of %d.\n", name, output[0].HighestBidder, output[0].HighestBid)

			fmt.Println("Auction has ended.")
			fmt.Printf("The highest bidder was %s with a bid of %d.\n", output[0].HighestBidder, output[0].HighestBid)
			wait <- true
			break
		}
	}

	<-wait
	elapsed := time.Since(start)
	log.Printf("Client: %v... Program done. Time taken: %s\n", name, elapsed.String())
	fmt.Println("Program done.")
	fmt.Printf("Time taken: %s\n", elapsed.String())
}

func ongoing(output []*proto.Outcome) bool {

	for _, outcome := range output {
		if outcome.Status == "Auction Ended" {
			return false
		}
	}

	return true
}

func prompt(stop chan bool) {
	inputchannel := make(chan string)

	go input(inputchannel, stop)

	for {
		select {
		case <-stop:
			stop <- true
			return
		case input := <-inputchannel:
			switch {
			case strings.HasPrefix(input, "bid"):
				var acks []*proto.Ack
				for _, auctionserver := range auctionserverconnections {
					bid(auctionserver, input, &acks)
				}
				if len(acks) > 0 {
					printSpaces()
					log.Printf("Client: %v... Bid with %s was: %s\n", name, input, acks[0].Outcome)
					fmt.Printf("Bid with %s was: %s\n", input, acks[0].Outcome)
				}
			case input == "result":
				var outcomes []*proto.Outcome
				for _, auctionserver := range auctionserverconnections {
					result(auctionserver, &outcomes)
				}
				if len(outcomes) > 0 {
					printSpaces()

					tense := "is"
					if outcomes[0].Status == "Auction Ended" {
						tense = "was"
					}
					log.Printf("Client: %v... !!!%v!!! The highest bid %v %d by %s.\n", name, outcomes[0].Status, tense, outcomes[0].HighestBid, outcomes[0].HighestBidder)
					fmt.Printf("!!!%v!!!\nThe highest bid %v %d by %s.\n", outcomes[0].Status, tense, outcomes[0].HighestBid, outcomes[0].HighestBidder)
				}
			}
		}
	}
}

func input(inputchannel chan string, stop chan bool) {
	for {
		select {
		case <-stop:
			stop <- true
			return
		default:
			reader := bufio.NewReader(os.Stdin)
			fmt.Println("To bid type 'Bid <amount>' and press enter.")
			fmt.Println("To get the status of the auction type 'Result' and press enter.")
			input, err := reader.ReadString('\n')
			if err != nil {
				log.Fatalln(err)
			}
			input = strings.TrimSpace(input)
			input = strings.ToLower(input)
			inputchannel <- input
			time.Sleep(1 * time.Second)
		}
	}
}

func bid(client proto.AuctionServiceClient, bid string, acks *[]*proto.Ack) {
	amountstr := strings.Split(bid, " ")[1]
	amountint, err := strconv.ParseInt(amountstr, 10, 64)
	if err != nil {
		log.Fatalln(err)
	}
	result, err := client.Bid(context.Background(), &proto.Amount{Bid: amountint, Bidder: name, BidTime: time.Now().Format(time.RFC850)})
	if err != nil {
		log.Fatalln(err)

	}
	*acks = append(*acks, result)
}

func result(client proto.AuctionServiceClient, outcomes *[]*proto.Outcome) {
	result, err := client.Result(context.Background(), &proto.Empty{})
	if err != nil {
		log.Fatalln(err)
	}
	*outcomes = append(*outcomes, result)
}

func printSpaces() {
	for i := 0; i < 100; i++ {
		fmt.Println()
	}
}

func connectToServers() {
	for _, server := range auctionservers {
		hostname := "localhost:" + server
		conn, err := grpc.NewClient(hostname, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatal(err.Error())
		}
		client := proto.NewAuctionServiceClient(conn)
		auctionserverconnections = append(auctionserverconnections, client)
	}
}

func getName() string {
	var name string
	var err error

	if len(os.Args) > 1 {

		name = os.Args[1]

	} else {
		fmt.Println("Enter your name:")
		reader := bufio.NewReader(os.Stdin)
		name, err = reader.ReadString('\n')
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	name = strings.TrimSpace(name)
	return name
}

func removePort(s []string, i int) []string {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
	//fmt.Println(append(s[:i], s[i+1:]...))
	//return append(s[:i], s[i+1:]...)
}

func removeConn(s []proto.AuctionServiceClient, i int) []proto.AuctionServiceClient {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
	//fmt.Println(append(s[:i], s[i+1:]...))
	//return append(s[:i], s[i+1:]...)
}
