package peer

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/libp2p/go-libp2p-core/peer"
	gorpc "github.com/libp2p/go-libp2p-gorpc"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/patrickmn/go-cache"
)

// A Person is made up of their owned peers and the relationships those peers have with other People
// A Person's peers do not have individual relationships with each other, they all write to the same Drama
type Person struct {
	ID    string
	Peers map[string]Drama // key is PeerID
}

type Self struct {
	Person
	PrivateKey crypto.PrivKey
	TLD        Drama // top-level drama that all peers resolve to
}

// TODO this should be it's own service? Wasn't able to get 2 services to work at the same time
// type PersonService struct {
// 	Cache *cache.Cache
// 	Peer  *Peer
// }

type CoordinateOTPArgs struct {
	CallingPeerID      string
	CoordinatingPeerID string
	OTP                string
}

// CoordinateOTP ingests an OTP from another peer you own who is coordinating a handshake
// ******** TODO these need to be authenticated
func (s *HandshakeService) CoordinateOTP(ctx context.Context, args CoordinateOTPArgs, reply *bool) error {
	// Time the process to record in the drama
	start := time.Now()

	// Receive an OTP and store in memory
	fmt.Println("Received an OTP:", args.OTP)
	s.Cache.Set(args.CallingPeerID, args.OTP, cache.DefaultExpiration)

	// TODO Coordinate the Self blockchain
	*reply = true
	fmt.Println("Process took:", time.Since(start))

	return nil
}

// func (p *Peer) listenForCoordination() {
// 	rpcHost := gorpc.NewServer(p.Host, HandshakeProtocolID)

// 	// Create a cache with a default expiration time of 5 minutes, and which
// 	// purges expired items every 10 minutes
// 	c := cache.New(5*time.Minute, 10*time.Minute)
// 	svc := HandshakeService{
// 		Cache: c,
// 		Peer:  p,
// 	}
// 	err := rpcHost.Register(&svc)
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println("Listening for coordination")
// }

func (p *Person) initiateCoordination(host host.Host, otp string) {
	rpcClient := gorpc.NewClient(host, HandshakeProtocolID)
	for pid, _ := range p.Peers {
		fmt.Println("Coordinating with peer:", pid)
		dest, err := peer.Decode(pid)
		if err != nil {
			log.Fatal(err)
		}

		args := CoordinateOTPArgs{
			CallingPeerID:      "TODO", // Okkk so is this the original request of the OTP or the Host that received the request that is now coordinating?
			CoordinatingPeerID: host.ID().String(),
		}
		var reply bool
		err = rpcClient.Call(dest, "HandshakeService", "CoordinateOTP", args, &reply)
		if err != nil {
			log.Fatal(err)
		}

		if reply {
			fmt.Println("Ok coordinated!")
		}
	}
}
