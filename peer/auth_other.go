// Authenticates OTHERS ---> essentially the server side
package peer

import "errors"

// This would be the (server) handler that would be hit on a queryPerson request
func authenticatePerson(p Person, app string) (int, error) {
	// First, make sure that all our peers got it. TODO: what happens if it times-out and junk?
	if !requestSentToAllPeers() {
		// TODO blacklist requester because it was invalid?
		// Definitely trigger an event or alert: https://opensource.com/article/19/10/event-driven-security
		return 0, errors.New("You suck")
	}

	// Ok so all peers are coordinating now

	return 100, nil
}

func requestSentToAllPeers() bool {
	return true
}
