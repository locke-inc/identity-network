package peer

import (
	"context"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/multiformats/go-multiaddr"
)

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

	// Use RPC to talk to the other peer
	// TODO how does rpc need to be authenticated?
	InitHandshake(p.Host, peerID)

}
