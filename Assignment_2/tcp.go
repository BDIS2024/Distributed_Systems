package main

import (
	"fmt"
	"time"
)

type TCPPackage struct {
	message string
	seq     int
	ack     int
}

type ChanSetup struct {
	in  chan TCPPackage
	out chan TCPPackage
}

func hasMessage(channel ChanSetup) bool {
	return len(channel.in) >= 1
}

func main() {
	channel1 := make(chan TCPPackage)
	channel2 := make(chan TCPPackage)
	done := make(chan int)

	go client(ChanSetup{channel1, channel2})
	go server(ChanSetup{channel2, channel1}, done)
	fmt.Printf("%s - Client and server spawned\n", time.Now().Format(time.StampNano))

	<-done
}

// Client methods

func client(channel ChanSetup) {
	var pack = TCPPackage{"request", 0, 0}
	clientSendSYN(channel, pack)
	if clientRecieveSYN(channel, &pack) { // if this method fails we end the handshake (and maybe try again)
		return
	}

	clientSendACK(channel, pack)
}

func clientSendSYN(channel ChanSetup, pack TCPPackage) { // We send SYN
	channel.out <- pack

	fmt.Printf("%s - Client sent SYN package: %v\n", time.Now().Format(time.StampNano), pack)
}

func clientRecieveSYN(channel ChanSetup, pack *TCPPackage) bool {
	var expectedAck = pack.seq + 1
	var failure = true

	select {
	case <-time.After(5 * time.Second): // Timeout after 5 seconds
		fmt.Printf("%s - Client timed out expecting package.\n", time.Now().Format(time.StampNano))

	case message := <-channel.in: // We recieve message
		if message.ack != expectedAck { // test if ack nr doesnt fit
			fmt.Printf("%s - Client recieved message with wrong ack nr. Expected: %d Recieved: %d\n", time.Now().String(), expectedAck, message.ack)
		} else {
			fmt.Printf("%s - Client recieved message: %v\n", time.Now().Format(time.StampNano), message)
			failure = false
		}
	}
	return failure
}

func clientSendACK(channel ChanSetup, pack TCPPackage) {
	pack.ack = pack.seq + 1

	channel.out <- pack

	fmt.Printf("%s - Client sent ACK package: %v\n", time.Now().Format(time.StampNano), pack)
}

// Server methods

func server(channel ChanSetup, done chan int) {
	message := <-channel.in
	fmt.Printf("%s - Server retrieved package with sequence: %d\n", time.Now().Format(time.StampNano), message.seq)
	message.seq += 1
	message.ack = 1

	channel.out <- message
	fmt.Printf("%s - Server sent back response with seq: %d and ack: %d\n", time.Now().Format(time.StampNano), message.seq, message.ack)

	response := <-channel.in
	if response.seq == 1 {
		channel.out <- response
	} else {
		fmt.Printf("%s - Server retrieved with sequence: %d, ack: %d and message: %s\n", time.Now().Format(time.StampNano), response.seq, response.ack, response.message)
	}
	done <- 1
}

/*
1)[Easy] Implement the TCP/IP Handshake using threads. This is not realistic (since the protocol should run across a network)
but your implementation needs to show that you have a good understanding of the protocol.
*/
