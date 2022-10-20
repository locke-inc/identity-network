package gateway

import (
	"fmt"
)

/*
 1. Accept incoming API requests (basically, allows non-peers to query peers).
    Later, perhaps we can allow servers to become peers themselves at will?
 2. Choose n nodes randomly to send incoming API requests, for example:
    Authenticate(person ID)
 3. Bootstrap peers. Bootstraps lookout nodes at same time
 4. Validate peer’s NameSystem and correct errors or detect attacks
*/
type Gateway struct {
	ID string
}

func (g *Gateway) Authenticate(peerID string) {
	g.ChooseNodes()
	// Query network
}

func (g *Gateway) ChooseNodes() {
	// Simply generates N amount of random addresses to query the network
	// The Gateway must somehow know how to randomly select ONLY occupied addresses
}

func (g *Gateway) Init() {
	fmt.Println("Initializing a new gateway node...")
}

func (g *Gateway) Stop() {
	fmt.Println("Stopping gateway node: ", g.ID)
}
