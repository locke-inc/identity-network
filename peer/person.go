package peer

type Person struct {
	ID           string
	Relationship Drama

	Peers map[string]Drama // key is PeerID
}
