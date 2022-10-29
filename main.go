package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/golang/protobuf/proto"
	"github.com/locke-inc/identity-network/peer"
	"github.com/locke-inc/identity-network/peer/pb"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

// The initial handshake
// type PeerHandshake struct {
// 	PeerID    string
// 	PublicKey string
// }

func main() {
	// Init a new peer ID
	peer := peer.New()
	fmt.Println("New peer initialized. Peer ID:", peer.PeerID)

	portPtr := flag.String("port", "3333", "port for TCP listener")
	flag.Parse()

	// Start the TCP server
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+*portPtr)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}

	// Close the listener when the application closes.
	defer l.Close()

	fmt.Println("Listening on " + CONN_HOST + ":" + *portPtr)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}

		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

	ph := &pb.PeerHandshake{}
	if err := proto.Unmarshal(buf, ph); err != nil {
		log.Fatalln("Failed to parse address book:", err)
		return
	}
	fmt.Println("Received peer handshake:", &ph)

	// Send a response back to person contacting us.
	conn.Write([]byte("Message received."))

	// Close the connection when you're done with it.
	conn.Close()
}

// func cmd(buf []byte) {
// 	switch(cmd) {
// 	case: ""
// 	}
// }
