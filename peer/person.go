package peer

type Person struct {
	ID           string
	Relationship Drama
	Peers        map[string]Drama // key is PeerID
}

func personIsTrusted() bool {
	// Get person from store
	// Send request to each peer with drama
	// If drama matches all peers, they all responds with an affirmative

	return true
}
