## TODO
I'm thinking that it makes sense to break down protocols into Events or Actions similar to how NGRX does it. I like the syntax of defining everything that happen, and the defining how it does happen.

Especially since we are dealing with Hosts -- devices that are both a client and a server. Both sends and receives data. So it can get sorta messy quickly. Would like to figure out an elegant way of defining the events that are sent and the events that are received.

FOR NOW using serveFunctionName and receiveFunctionName to clarify a lil

## Personal nodes (peers):
Peers are used for locating other peers using the DHT, and are also used to ask questions (zero-knowledge proofs) to verify that a peer is who they claim to be. A lookout node is just a personal node hosted by Locke with a few additional algorithms to provide assurances.

1. Each peer bootstraps itself with a Gateway node
2. Each peer has a DHT list of people it knows
3. Each peer has a NameSystem to enable human readable usernames
4. Each peer can handshake with other peers using a private key made by Crystals Kyber for e2e encryption between them. This looks the same for lookout nodes as it does for personal nodes.
5. Each peer participates in community auth (peers coming to a consensus of how likely a given peer is to be who they claim to be). MAYBE?: Each auth attempt has a blockchain (essentially a merkle tree) as the message structure to create immutable receipts.

## Community store
Using boltDB, each person has their own bucket. Inside that bucket is a list of their peerIDs as the key, and their blockchain message history as the value

## Terminology
Peer vs node - peer is all layer 0 devices, gateway nodes are layer 2

