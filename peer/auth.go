package peer

import (
	"context"

	"github.com/libp2p/go-libp2p-core/peer"
)

func authenticate(ctx context.Context, p *Peer, application string) {
	// Application string defines what the auth request context is
	// Start by sending a request to all people in your store
	people, err := p.getAllPeople()
	if err != nil {
		panic(err)
	}

	var keys []string
	for i, person := range people {
		keys[i] = queryPerson(ctx, p, person, application)
	}

}

// TODO generic af name
func queryPerson(ctx context.Context, p *Peer, person Person, app string) string {
	for peerID, _ := range person.Peers {
		pid, err := peer.Decode(peerID)
		if err != nil {
			panic(err)
		}

		// TODO Send request to peerID
		str, err := p.Host.NewStream(ctx, pid, "/locke/1.0.0")
		handleStream(str)
	}

	return "this would be a key shard"
}
