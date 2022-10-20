package peer

import (
	"fmt"

	"github.com/ipfs/go-namesys"
	dht "github.com/libp2p/go-libp2p-kad-dht"
)

/*
 1. Each peer bootstraps itself with a Gateway node
 2. Each peer has a DHT list of people it knows
 3. Each peer has a NameSystem to enable human readable usernames
 4. Each peer can handshake with other peers using a private key made by
    Crystals Kyber for e2e encryption between them. This looks the same for
    lookout nodes as it does for personal nodes.
 5. Each peer participates in community auth (peers coming to a consensus
    of how likely a given peer is to be who they claim to be). MAYBE?: Each
    auth attempt has a blockchain (essentially a merkle tree) as the message
    structure to create immutable receipts.
*/
type Peer struct {
	Identity
	DHT     *dht.IpfsDHT
	Namesys namesys.NameSystem
}

type Identity struct {
	PeerID  string
	PrivKey string `json:",omitempty"`
}

func (p *Peer) Bootstrap() {
	fmt.Println("Initializing new peer")
}

func (p *Peer) Stop() {
	fmt.Println("Stopping peer: ", p.Identity.PeerID)
}
