package peer

import (
	"context"
	"time"

	"github.com/libp2p/go-libp2p-core/peer"
)

type Transaction struct {
	Requester   string // peerID that initiated the transaction
	RequestType string // "auth" for example

	Responder   string // peerID that responded
	Result      int    // for now let's say it's between 0 and 100%
	ProcessTime time.Duration
}

func transaction(ctx context.Context, p *Peer) {
	// Start by sending a request to all people
	people, err := p.getAllPeople()
	if err != nil {
		panic(err)
	}

	for _, person := range people {
		queryPerson(ctx, p, person)
	}
}

// TODO generic af name
func queryPerson(ctx context.Context, p *Peer, person Person) {
	for peerID, _ := range person.Peers {
		pid, err := peer.Decode(peerID)
		if err != nil {
			panic(err)
		}

		// TODO Send request to peerID
		str, err := p.Host.NewStream(ctx, pid, "/locke/1.0.0")
		handleStream(str)
	}
}
