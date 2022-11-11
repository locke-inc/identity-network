package peer

import (
	"fmt"
	"time"

	gorpc "github.com/libp2p/go-libp2p-gorpc"
	"github.com/patrickmn/go-cache"
)

// Handshake protocol
func (p *Peer) listenForHandshake() {
	rpcHost := gorpc.NewServer(p.Host, HandshakeProtocolID)

	// Create a cache with a default expiration time of 5 minutes, and which
	// purges expired items every 10 minutes
	c := cache.New(5*time.Minute, 10*time.Minute)
	svc := HandshakeService{
		Cache: c,
		Peer:  p,
	}
	err := rpcHost.Register(&svc)
	if err != nil {
		panic(err)
	}

	fmt.Println("\nListening for handshakes")
}

// Authenticate protocol
func (p *Peer) listenForAuthenticate() {
	rpcHost := gorpc.NewServer(p.Host, HandshakeProtocolID)

	// Create a cache with a default expiration time of 5 minutes, and which
	// purges expired items every 10 minutes
	c := cache.New(5*time.Minute, 10*time.Minute)
	svc := HandshakeService{
		Cache: c,
		Peer:  p,
	}
	err := rpcHost.Register(&svc)
	if err != nil {
		panic(err)
	}

	fmt.Println("\nListening for handshakes")
}
