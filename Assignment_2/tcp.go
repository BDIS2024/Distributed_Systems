package main

/*

	Authored by: Markus Sværke Staael 	(msvs@itu.dk)
	Authored by: Patrick Shen			(pash@itu.dk)
	Authored by: Nicky Chengde Ye		(niye@itu.dk)

	This sollution is a demonstration of a TCP handshake.
	In this "world" we assume there is only one client and server that will communicate with each other.
	The client or server will not retry if any errors occur (there will not occur errors) but will check for them to demonstrate knowledge of what is going to happen.
	The client or server will not have timeouts for sending and recieving messages.

*/

import (
	"fmt"
	"math/rand/v2"
	"time"
)

type TCPHeader struct {
	sequence_number       int
	acknowledgment_number int
	ACK                   bool
	SYN                   bool
}

type ChanSetup struct {
	in  chan TCPHeader
	out chan TCPHeader
}

func hasMessage(channel ChanSetup) bool {
	return len(channel.in) >= 1
}

func main() {
	channel1 := make(chan TCPHeader)
	channel2 := make(chan TCPHeader)

	go client(ChanSetup{channel1, channel2})
	go server(ChanSetup{channel2, channel1})

	for {
	}
}

// CLIENT
func client(channel ChanSetup) {
	fmt.Printf("%s - Client routine started.\n", time.Now().Format(time.StampNano))

	//
	// STEP 1: Create our first connection request.
	//

	var client_ISN = rand.IntN(100)
	var syn_package = TCPHeader{
		sequence_number:       client_ISN, // Client random sequence number
		acknowledgment_number: 0,
		ACK:                   false,
		SYN:                   true, // We want to estabish connection
	}
	channel.out <- syn_package
	fmt.Printf("%s - Client sent SYN package: %v\n", time.Now().Format(time.StampNano), syn_package)

	//
	// STEP 3: Recieve the SYN-ACK from the server
	//

	SYN_ACK_message_recieve := <-channel.in
	fmt.Printf("%s - Client recieved SYN-ACK package: %v\n", time.Now().Format(time.StampNano), SYN_ACK_message_recieve)

	// STEP 3.1: Check if all the recieved data is correct

	if !(SYN_ACK_message_recieve.SYN == true) && // Assert that SYN flag is true
		!(SYN_ACK_message_recieve.ACK == true) && // Assert that ACK flag is true
		!(SYN_ACK_message_recieve.acknowledgment_number == client_ISN+1) { // Assert that acknowledgement number is as expected
		clientEndDemonstration()
		return
	}
	var server_ISN = SYN_ACK_message_recieve.sequence_number

	// STEP 3.2 return ACK package
	var ack_package = TCPHeader{
		sequence_number:       client_ISN + 1, // Client random sequence number
		acknowledgment_number: server_ISN + 1,
		ACK:                   true,
		SYN:                   false, // We want to estabish connection
	}

	fmt.Printf("%s - Client sent ACK package: %v\n", time.Now().Format(time.StampNano), ack_package)

	// Connection has been established
}

// SERVER
func server(channel ChanSetup) {
	fmt.Printf("%s - Server listening.\n", time.Now().Format(time.StampNano))

	//
	// STEP 2: RECIEVE SYN and sent SYN-ACK
	//

	SYN_message_recieve := <-channel.in
	fmt.Printf("%s - Server retrieved package: %v\n", time.Now().Format(time.StampNano), SYN_message_recieve)

	// 2.1 Test to see if client want to establish connection (the expected message for this simple demonstration)
	if SYN_message_recieve.SYN == false {
		serverEndDemonstration(SYN_message_recieve)
		return
	}

	// 2.2 Create a the returning SYN-ACK message
	var server_ISN = rand.IntN(100)                      // The servers Sequence number
	var client_ISN = SYN_message_recieve.sequence_number // The client sequence number (what sequence number we expect from them)
	var syn_ack_package = TCPHeader{
		sequence_number:       server_ISN,
		acknowledgment_number: client_ISN + 1,
		ACK:                   true, // We aknowledge your sequence number
		SYN:                   true, // We also want to estabish connection
	}

	// 2.3 Return a reply
	channel.out <- syn_ack_package
	fmt.Printf("%s - Server sent back SYN-ACK response: %v\n", time.Now().Format(time.StampNano), syn_ack_package)

	//
	// (Any steps from here will not be demonstrated)
	//
	<-channel.in
}

func clientEndDemonstration() {
	fmt.Printf("%s - Client recieved unexpected package terminating demonstration.\n", time.Now().Format(time.StampNano))
}

func serverEndDemonstration(arg TCPHeader) {
	fmt.Printf("%s - Server recieved unexpected package (NO SYN FLAG) terminating demonstration.\n", time.Now().Format(time.StampNano))
}
