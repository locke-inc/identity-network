package main

import (
	"fmt"

	"github.com/locke-inc/identity-network/peer"
)

func main() {
	peer := peer.New()
	fmt.Println("A new peer has been created. PeerID is:", peer.Identity.PeerID)
}
