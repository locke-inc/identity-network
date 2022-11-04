# identity-network

The main functionality we're building is resolving peers to people.
Ultimately the goal is for each person to have a lot of devices in their layer 0.

Build for linux:
env GOOS=linux GOARCH=amd64 go build

## TODO
Eddilithium3 is *almost* there but it needs to be implemented into the OpenSSL protocol as a key option, should be plug and play after that
    Maybe bring the entire QUIC protocol into this repo rather than using libp2p's

