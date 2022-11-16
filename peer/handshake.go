package peer

import (
	"bytes"
	"context"
	cryptorand "crypto/rand"
	"encoding/gob"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	gorpc "github.com/libp2p/go-libp2p-gorpc"
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

// (1) Ask to start a new relationship
type StartRelationshipArgs struct {
	CallingPeerID string
}

type StartRelationshipResp struct {
	ReadyForAuth bool
}

// (2) Authorize if the relationship should happen
type AuthorizeRelationshipArgs struct {
	CallingPeerID string
	Them          Person
	OTP           string
}

type AuthorizeRelationshipResp struct {
	Them   Person
	TLD    Drama // top-level drama
	SymKey []byte
}

// (3) Settle the relationship
type SettleRelationshipArgs struct {
	CallingPeerID string
	Them          Person
	Drama
}

type SettleRelationshipResp struct {
	Success bool
}

// StartRelationship is a peer-to-peer function, with this being the receiver peer.
// The receiver peer handles the coordination of everything
func (s *HandshakeService) StartRelationship(ctx context.Context, args StartRelationshipArgs, resp *StartRelationshipResp) error {
	// Time the process to record in the drama
	start := time.Now()

	// Generate an OTP
	fmt.Println("Generating an OTP")
	otp, err := generateOTP(6)
	if err != nil {
		return err
	}

	// Display OTP for manual sharing
	fmt.Println("\n> OTP: ", otp)

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
	drama := CreateDrama()

	// Add an unencrypted block since this peer is not established as trusted yet and this information should be public
	drama.addUnencryptedBlock(t)

	// Store
	err = s.Peer.addPeer(TempPeerBucket, args.CallingPeerID, &drama)
	if err != nil {
		return err
	}

	resp.ReadyForAuth = true
	return nil
}

func (s *HandshakeService) AuthorizeRelationship(ctx context.Context, args AuthorizeRelationshipArgs, resp *AuthorizeRelationshipResp) error {
	// Time the process to record in the drama
	start := time.Now()

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
	// fmt.Println("Got stored drama", drama)
	if !drama.isValid() {
		fmt.Println("Drama is NOT valid, abort")
		return errors.New("Drama is invalid")
	}

	// Verify the OTP is correct
	var cachedOTP string
	x, found := s.Cache.Get(args.CallingPeerID)
	if !found {
		return errors.New("No cached OTP to compare")
	}

	cachedOTP = x.(string)

	if strings.TrimRight(args.OTP, "\n") != cachedOTP {
		// TODO
		// OTP did not match, record this transaction as failed and add a maximum number of calls before this peer is blacklisted
		return errors.New("OTP did not match")
	}

	fmt.Println("\nHandshake success.")

	// OTP matched, generate sym key for this relationship and add person
	symKey := make([]byte, chacha20poly1305.KeySize)
	if _, err := cryptorand.Read(symKey); err != nil {
		panic(err)
	}

	s.Peer.addPerson(&args.Them, &drama, &symKey)

	// Remove the temp peer
	s.Peer.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(TempPeerBucket))
		b.Delete([]byte(Prefix_Peer + args.CallingPeerID))
		return nil
	})

	t := Transaction{
		Requester:   args.CallingPeerID,
		RequestType: "handshake",
		Responder:   s.Peer.Self.ID,
		Result:      99, // 99 represents a successful OTP auth <---- this is a little cheeky; it's not 100 since we're never 100% sure of anything...
		Application: "handshake success",
		ProcessTime: time.Since(start),
	}
	drama.addBlock(t, symKey)

	// Them is you to them!
	resp.Them = s.Peer.Self.Person
	resp.TLD = drama
	resp.SymKey = symKey
	return nil
}

func (s *HandshakeService) SettleRelationship(ctx context.Context, args SettleRelationshipArgs, resp *SettleRelationshipResp) error {
	if !args.isValid() {
		fmt.Println("Drama is NOT valid, abort")
		resp.Success = false
		return errors.New("Drama is invalid")
	}

	err := s.Peer.updateDrama(args.Them.ID, args.CallingPeerID, &args.Drama)
	if err != nil {
		resp.Success = false
		return err
	}

	resp.Success = true
	return nil
}

// Handshake protocol
func (p *Peer) listenForHandshake() {
	rpcHost := gorpc.NewServer(p.Host, HandshakeProtocolID)

	// Create a cache with a default expiration time of 5 minutes, and which
	// purges expired items every 10 minutes
	c := cache.New(5*time.Minute, 10*time.Minute)
	svc := HandshakeService{
		Cache: c,
		Peer:  p,
	}
	err := rpcHost.Register(&svc)
	if err != nil {
		panic(err)
	}

	fmt.Println("\nListening for handshakes")
}

// TODO https://github.com/libp2p/specs/blob/master/discovery/mdns.md - Allows peers on same network to discover each other easily
