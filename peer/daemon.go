package peer

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/multiformats/go-multiaddr"
)

func handleStream(s network.Stream) {
	log.Println("Got a new stream!")

	// Create a buffer stream for non blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	go readData(rw)
	go writeData(rw)
}

func readData(rw *bufio.ReadWriter) {
	for {
		str, _ := rw.ReadString('\n')

		if str == "" {
			return
		}
		if str != "\n" {
			// Green console colour: 	\x1b[32m
			// Reset console colour: 	\x1b[0m
			fmt.Printf("\x1b[32m%s\x1b[0m> ", str)
		}
	}
}

func writeData(rw *bufio.ReadWriter) {
	stdReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			log.Println(err)
			return
		}

		rw.WriteString(fmt.Sprintf("%s\n", sendData))
		rw.Flush()
	}
}

func connect(ctx context.Context, p *Peer, destination string, pid string) {
	peerID, err := peer.Decode(pid)
	if err != nil {
		panic(err)
	}

	addr, err := multiaddr.NewMultiaddr(destination)
	if err != nil {
		panic(err)
	}

	// Add the destination's peer multiaddress in the peerstore.
	// This will be used during connection and stream creation by libp2p.
	var maddr []multiaddr.Multiaddr
	p.Host.Peerstore().AddAddrs(peerID, append(maddr, addr), peerstore.PermanentAddrTTL)

	// Start a stream with the destination.
	// Multiaddress of the destination peer is fetched from the peerstore using 'peerId'.
	// str, err := p.Host.NewStream(ctx, peerID, HandshakeProtocolID)
	// if err != nil {
	// 	log.Println(err)
	// 	panic(err)
	// }

	// log.Println("Established connection to destination")

	// // Send a handshake
	// p.identifySelf(str)

	initiateHandshake(p.Host, peerID)

}
