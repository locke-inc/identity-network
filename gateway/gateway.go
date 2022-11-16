package gateway

import (
	"context"
	"fmt"

	"github.com/locke-inc/identity-network/peer"
)

type Gateway struct {
	Peer peer.Peer
}

func (g *Gateway) New() {
	fmt.Println("Creating new gateway")
	p := peer.Peer{}
	p.New()
	g.Peer = p

	// Listen for gateway specific calls
	g.listenForCentralAuth()

	select {}
}

func (g *Gateway) test() {
	// Let's test the RPC
	svc := CentralAuthService{
		Peer: &g.Peer,
	}

	args := CentralAuthArgs{
		Personame:    "tester100",
		PasswordHash: "xKfDkCyBLDl5YoxLPtWoPwqW4F1eNHAxXskb/E+zOgo=",
	}
	var resp CentralAuthResp
	err := svc.CentralAuth(context.Background(), args, &resp)
	if err != nil {
		panic(err)
	}

	fmt.Println("Resp:", resp)
}
