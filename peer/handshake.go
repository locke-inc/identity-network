package peer

import "fmt"

// TODO https://github.com/libp2p/specs/blob/master/discovery/mdns.md
// Allows peers on same network to discover each other easily
func (p *Peer) handshake(pid string) {
	// TODO Check if this peer is already in store
	// TODO identify owner of peer somehow
	// Add peer to community store
	var peers []string
	err := p.addPerson("connor", append(peers, pid))
	if err != nil {
		panic(err)
	}

	fmt.Println("Added " + pid + " to local store")

	// Should generate a shared sym key here? Prolly

	// Broadcast this to all other owned peers so they can handshake
	// If they don't ALL return a "success" msg then it failed
}
