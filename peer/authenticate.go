/*
https://en.wikipedia.org/wiki/Randomized_algorithm

A randomized algorithm is an algorithm that employs a degree of randomness as part of its logic or procedure.

One has to distinguish between algorithms that use the random input so that they always terminate with the correct answer,
but where the expected running time is finite (Las Vegas algorithms, for example Quicksort[1]), and algorithms which have a
chance of producing an incorrect result (Monte Carlo algorithms, for example the Monte Carlo algorithm for the MFAS problem[2])
or fail to produce a result either by signaling a failure or failing to terminate.


***In some cases, probabilistic algorithms are the only practical means of solving a problem.[3]


In common practice, randomized algorithms are approximated using a pseudorandom number generator in place of a true source of
random bits; such an implementation may deviate from the expected theoretical behavior and mathematical guarantees which may


**** depend on the existence of an ideal true random number generator. <------ 	theory: passing around messages in a chaotic
																				network generates truly quantum effects
*/

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

type AuthPersonArgs struct {
	Caller string
	Key    string // what key do they need to send to prove that it's a legit request?
}

func InitAuthentication(host host.Host, dest peer.ID) {
	fmt.Println("Initiating authentication...")
	rpcClient := gorpc.NewClient(host, AuthProtocolID)

	// First call
	args := StartRelationshipArgs{
		CallingPeerID: host.ID().String(),
	}
	var reply Drama
	err := rpcClient.Call(dest, "AuthService", "AuthenticatePerson", args, &reply)
	if err != nil {
		log.Fatal(err)
	}
}

// Let's say that this is the function that the gateways call directly
// So it also includes locating the person
func (s *AuthService) AuthenticatePerson(ctx context.Context, args AuthPersonArgs, reply *Person) error {
	// Locate person through the DHT somehow
	return nil
}

func authenticateSelf(ctx context.Context, p *Peer, application string) {
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
		_, err := peer.Decode(peerID)
		if err != nil {
			panic(err)
		}

		// TODO Send request to peerID
		// str, err := p.Host.NewStream(ctx, pid, "/locke/1.0.0")
		// handleStream(str)
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
