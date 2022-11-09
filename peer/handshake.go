package peer

import (
	"bufio"
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"log"

	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	gorpc "github.com/libp2p/go-libp2p-gorpc"
	"github.com/libp2p/go-libp2p/core/host"
)

const (
	HandshakeProtocolID = "/locke/handshake"
	// IdentifyYourselfEndpoint      = HandshakeProtocolID + "/id_self/1.0.0"
	// StartRelationshipEndpoint = HandshakeProtocolID + "/start_relationship/1.0.0"
)

type HandshakeArgs struct {
	Key []byte
}
type HandshakeReply struct {
	Who Person
}
type HandshakeService struct {
	Person
}

func (me *HandshakeService) IdentifyYourself(ctx context.Context, argType HandshakeArgs, replyType *HandshakeReply) error {
	log.Println("Received a Ping call")
	replyType.Who = me.Person
	return nil
}

func (t *HandshakeService) StartRelationship(ctx context.Context, argType HandshakeArgs, replyType *HandshakeReply) error {
	// log.Println("Received a Ping call")
	// replyType.Data = argType.Data
	return nil
}

func initiateHandshake(host host.Host, dest peer.ID) {
	rpcClient := gorpc.NewClient(host, HandshakeProtocolID)

	var reply HandshakeReply
	var args HandshakeArgs

	// Auth would happen here
	b := make([]byte, 32)
	args.Key = b

	err := rpcClient.Call(dest, "HandshakeService", "IdentifyYourself", args, &reply)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Got the call back I guess:", reply.Who)
}

func (p *Peer) listenForHandshake() {
	rpcHost := gorpc.NewServer(p.Host, HandshakeProtocolID)

	svc := HandshakeService{
		p.Me,
	}
	err := rpcHost.Register(&svc)
	if err != nil {
		panic(err)
	}

	fmt.Println("Done")
}

// (1) A handshake starts by each peer identifying their owner (self)
func (p *Peer) identifySelf(s network.Stream) {
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	go sendSelf(rw, p.Me)
	go receiveThem(rw, p)
}

func sendSelf(stream *bufio.ReadWriter, self Person) {
	fmt.Println("Sending self...")

	// Encode self and send downstream
	var conn bytes.Buffer
	err := gob.NewEncoder(&conn).Encode(self)
	if err != nil {
		panic(err)
	}

	i, err := stream.Write(conn.Bytes())
	if i == 0 || err != nil {
		panic(err)
	}
	stream.Flush()
}

func receiveThem(stream *bufio.ReadWriter, p *Peer) {
	// Read from stream and decode gob
	var them Person
	if err := gob.NewDecoder(stream).Decode(&them); err != nil {
		panic(err)
	}

	// Store
	p.addNewPerson(them)
	fmt.Println("Received them:", them)
}

// (2) A handshake continues by the peers establishing a new drama
func startDrama(s network.Stream) {
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	go sendDrama(rw)
	go receiveDrama(rw)
}

func sendDrama(stream *bufio.ReadWriter) {
	fmt.Println("Sending drama")
	var conn bytes.Buffer
	var d = CreateDrama(0)
	err := gob.NewEncoder(&conn).Encode(d)
	if err != nil {
		panic(err)
	}

	i, err := stream.Write(conn.Bytes())
	if i == 0 || err != nil {
		panic(err)
	}
	stream.Flush()
}

func receiveDrama(stream *bufio.ReadWriter) {
	var drama Drama
	if err := gob.NewDecoder(stream).Decode(&drama); err != nil {
		panic(err)
	}
	fmt.Println("New drama received: check it:", drama)
}

// TODO https://github.com/libp2p/specs/blob/master/discovery/mdns.md - Allows peers on same network to discover each other easily
