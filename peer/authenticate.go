package peer

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"

	"github.com/libp2p/go-libp2p-core/peer"
	gorpc "github.com/libp2p/go-libp2p-gorpc"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/patrickmn/go-cache"
)

const (
	AuthProtocolID = "/locke/auth"
)

type AuthService struct {
	Cache *cache.Cache
	Peer  *Peer
}

func InitAuthentication(host host.Host, dest peer.ID) {
	fmt.Println("Initiating authentication...")
	rpcClient := gorpc.NewClient(host, AuthProtocolID)

	// First call
	args1 := StartRelationshipArgs{
		CallingPeerID: host.ID().String(),
	}
	var reply Drama
	err := rpcClient.Call(dest, "AuthService", "StartRelationship", args1, &reply)
	if err != nil {
		log.Fatal(err)
	}
}

// Let's say that this is the function that the gateways call directly
// So it also includes locating the person
func (s *AuthService) AuthenticatePerson(ctx context.Context, args VerifyOTPArgs, reply *Person) error {
	return nil
}

func authenticate(ctx context.Context, p *Peer, application string) {
	people, err := p.getAllPeople()
	if err != nil {
		panic(err)
	}

	// Send request to everyone in community
	var keys []string
	for i, person := range people {
		keys[i] = queryPerson(ctx, p, person, application)
	}

}

func queryPerson(ctx context.Context, p *Peer, person Person, app string) string {
	for peerID := range person.Peers {
		pid, err := peer.Decode(peerID)
		if err != nil {
			panic(err)
		}

		// TODO Send request to peerID
		str, err := p.Host.NewStream(ctx, pid, "/locke/1.0.0")
		handleStream(str)
	}

	return "this would be a key shard"
}

// Helpers
const otpChars = "1234567890"

func generateOTP(length int) (string, error) {
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	otpCharsLength := len(otpChars)
	for i := 0; i < length; i++ {
		buffer[i] = otpChars[int(buffer[i])%otpCharsLength]
	}

	return string(buffer), nil
}
