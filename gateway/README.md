## Gateway Nodes
Gateway nodes DO NOT have a list of all peers. Each peer has a list of the peers it knows, creating a distributed hash table. When an API request comes in, the gateway node sends the request to n random nodes. The nodes locate the requesting node and ask it to authenticate itself (zero-knowledge proofs).

1. Accept incoming API requests (basically, allows non-peers to query peers). Later, perhaps we can allow servers to become peers themselves at will?
2. Choose n nodes randomly to send incoming API requests, for example: Authenticate(person ID)
3. Bootstrap peers. Bootstraps lookout nodes at same time
4. Validate peer’s NameSystem and correct errors or detect attacks
