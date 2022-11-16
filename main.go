package main

import (
	"flag"

	"github.com/locke-inc/identity-network/gateway"
	"github.com/locke-inc/identity-network/peer"
)

func main() {
	// Get flags
	nodeType := flag.String("type", "peer", "What type of node do you want?")
	flag.Parse()

	switch *nodeType {
	case "peer":
		p := peer.Peer{}
		p.New()
	case "gateway":
		g := gateway.Gateway{}
		g.New()
	}

}
