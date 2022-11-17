package main

import (
	"flag"
	"fmt"

	"github.com/locke-inc/identity-network/gateway"
	"github.com/locke-inc/identity-network/peer"
)

func main() {
	// Get flags
	nodeType := flag.String("type", "peer", "What type of node do you want?")
	personName := flag.String("name", "", "What's your Locke username?")
	flag.Parse()

	if *personName == "" {
		fmt.Println("No username defined: \"-name\" flag was not defined")
		return
	}

	switch *nodeType {
	case "peer":
		p := peer.Peer{}
		p.Start(*personName)
	case "gateway":
		g := gateway.Gateway{}
		g.Start(*personName)
	}
}
