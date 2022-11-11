# identity-network

The main functionality we're building is resolving peers to people.
Ultimately the goal is for each person to have a lot of devices in their layer 0.

Build for linux:
env GOOS=linux GOARCH=amd64 go build

## TODO
1. Eddilithium3 is *almost* there but it needs to be implemented into the OpenSSL protocol as a key option, should be plug and play after that
    Maybe bring the entire QUIC protocol into this repo rather than using libp2p's

2. Onboarding peers to people. Their first node is a separate process, since future nodes need to know of their existence.

3. Peer coordination/syncing
    Super important in order to resolve people to peers and vice versa
    There are 2 ways a message can be sent to a Person:
        1. The message is sent to every peer the Person owns at the same time. The peers then validate that the message was received by all online peers.
        2. The message is sent to a single peer who then distributes the message and "coordinates" the response.

4. Auth from the Locke API, forwarded by gateway nodes.
    I think for now gateway nodes