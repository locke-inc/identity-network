package peer

import (
	"bufio"
	"bytes"
	"context"
	cryptorand "crypto/rand"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"github.com/libp2p/go-libp2p-core/peer"
	gorpc "github.com/libp2p/go-libp2p-gorpc"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/patrickmn/go-cache"
	"golang.org/x/crypto/chacha20poly1305"
)

const (
	HandshakeProtocolID = "/locke/handshake"
	TempPeerBucket      = "tmp"
)

type HandshakeService struct {
	Cache *cache.Cache
	Peer  *Peer
}

type StartRelationshipArgs struct {
	CallingPeerID string
}

type VerifyOTPArgs struct {
	CallingPeerID string
	OTP           string
}

type VerifyOTPResponse struct {
	Them   Person
	SymKey []byte
}

func InitHandshake(host host.Host, dest peer.ID) {
	fmt.Println("Initiating handshake...")
	rpcClient := gorpc.NewClient(host, HandshakeProtocolID)

	args1 := StartRelationshipArgs{
		CallingPeerID: host.ID().String(),
	}
	var reply Drama
	err := rpcClient.Call(dest, "HandshakeService", "StartRelationship", args1, &reply)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Ok, they're ready for me to send them an OTP. Here's the drama:\n", reply)

	// Input OTP
	stdReader := bufio.NewReader(os.Stdin)
	fmt.Print("> ")
	otp, err := stdReader.ReadString('\n')
	if err != nil {
		log.Println(err)
		return
	}

	args2 := VerifyOTPArgs{
		CallingPeerID: host.ID().String(),
		OTP:           otp,
	}
	var resp VerifyOTPResponse
	err = rpcClient.Call(dest, "HandshakeService", "VerifyOTP", args2, &resp)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Ok got them!", resp)

	// TODO now need to send them ourselves to finish the handshake
	// Store symkey and all that
}

// StartRelationship is a peer-to-peer function, with this being the receiver peer.
// The receiver peer handles the coordination of everything
func (s *HandshakeService) StartRelationship(ctx context.Context, args StartRelationshipArgs, reply *Drama) error {
	// Time the process to record in the drama
	start := time.Now()

	// Generate an OTP
	fmt.Println("Generating an OTP")
	otp, err := generateOTP(6)
	if err != nil {
		return err
	}

	// Display OTP for manual sharing
	fmt.Println(otp)

	// Cache OTP for later comparison
	s.Cache.Set(args.CallingPeerID, otp, cache.DefaultExpiration)

	// TODO Send the otp to your other peers so any of them can authenticate the start of this relationship
	// s.Peer.Self.initiateCoordination(s.Peer.Host, otp)

	// Create new transaction to record into drama
	t := Transaction{
		Requester:   args.CallingPeerID,
		RequestType: "handshake",
		Responder:   s.Peer.Self.ID,
		Result:      0,
		Application: "init handshake",
		ProcessTime: time.Since(start),
	}
	drama := CreateDrama(0)

	// Add an unencrypted block since this peer is not established as trusted yet and this information should be public
	drama.addUnencryptedBlock(t)
	fmt.Println("Drama made. Storing...", drama)
	err = s.Peer.addPeer(TempPeerBucket, args.CallingPeerID, &drama)
	if err != nil {
		return err
	}

	// Send back
	*reply = drama
	fmt.Println("Process took:", t.ProcessTime)
	return nil
}

func (s *HandshakeService) VerifyOTP(ctx context.Context, args VerifyOTPArgs, reply *VerifyOTPResponse) error {
	// First check this peer's history
	var drama Drama
	err := s.Peer.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(TempPeerBucket))
		v := b.Get([]byte(Prefix_Peer + args.CallingPeerID))
		// Decode drama from gob
		buf := bytes.NewBuffer(v)
		dec := gob.NewDecoder(buf)
		if err := dec.Decode(&drama); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	// TODO
	// Parse drama to ensure it's valid and peer can be trusted
	// For instance if it's blacklisted can be denied completely
	fmt.Println("Got drama from store:", drama)

	// Verify the OTP is correct
	log.Println("Verifying OTP, received:", args.CallingPeerID)
	var cachedOTP string
	x, found := s.Cache.Get(args.CallingPeerID)
	if !found {
		return errors.New("No cached OTP to compare")
	}

	cachedOTP = x.(string)

	if strings.TrimRight(args.OTP, "\n") != cachedOTP {
		return errors.New("OTP did not match")
	}

	// OTP matched, generate sym key for this relationship
	symKey := make([]byte, chacha20poly1305.KeySize)
	if _, err := cryptorand.Read(symKey); err != nil {
		panic(err)
	}

	// Cache symKey, will be added to persistent store once they respond with their identity
	s.Cache.Set(args.CallingPeerID, symKey, cache.DefaultExpiration)

	// Them is you to them!
	reply.Them = s.Peer.Self
	return nil
}

// TODO https://github.com/libp2p/specs/blob/master/discovery/mdns.md - Allows peers on same network to discover each other easily
