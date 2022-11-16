package lookout

import "github.com/locke-inc/identity-network/peer"

type Lookout struct {
	Peer peer.Peer
}

func LoadPeer() {
	// Starting a lookout by pulling in a private key from a peer_ID --> priv_key key store
	// DynamoDB??

	// Uh oh, can a lambda be invoked from an RPC call?
}
