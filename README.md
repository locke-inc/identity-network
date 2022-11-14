# identity-network

The main functionality we're building is defining a Person datatype for the Internet. A Person is defined by the devices they own and the relationships those peers have with other people's peers.

We resolving peers to people.

Ultimately the goal is for each person to have a lot of devices in their layer 0, Locke will sell a hosted peer and maybe a hardware device that plugs into a home router at some point.

## Run
go run main.go

Starts your own local peer. To connect to another peer run:
go run main.go -dest /ip4/$PeerPublicIP/udp/5533/quic -peer $PeerID

Replacing $PeerPublicIP with a ipv4 address and $PeerID with the peerID located at that address.

Of course we will eventually enable DHT routing by peer ID and you won't need to enter in their IP address, but for now it's fine.

## TODO
1. Eddilithium3 is *almost* there but it needs to be implemented into the OpenSSL protocol as a key option, should be plug and play after that
    Maybe bring the entire QUIC protocol into this repo rather than using libp2p's

2. Onboarding peers to people. Their first node is a separate process, since future nodes need to know of their existence. Initial auth should come from Locke API and maybe be forwarded from gateway nodes.

3. Peer coordination/syncing
    Super important in order to resolve people to peers and vice versa
    There are 2 ways a message can be sent to a Person:
        1. The message is sent to every peer the Person owns at the same time. The peers then validate that the message was received by all online peers.
        2. The message is sent to a single peer who then distributes the message and "coordinates" the response.