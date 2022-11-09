package peer

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/libp2p/go-libp2p-core/network"
)

type EventHandshake struct {
}

func (p *Peer) handleHandshake(s network.Stream) {
	// Create a buffer stream for non blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	go serveHandshake(rw, p.Me)
	go receiveHandshake(rw, p)
}

func serveHandshake(stream *bufio.ReadWriter, self Person) {
	fmt.Println("Serving a handshake!", self)
	sendSelf(stream, self)
	sendDrama(stream, self)
	// Create a new drama for this relationship
	// TODO how is it decided who starts the relationship? Whoever serves first right?

}

func sendSelf(stream *bufio.ReadWriter, self Person) {
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

func sendDrama(stream *bufio.ReadWriter, self Person) {
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

func receiveHandshake(stream *bufio.ReadWriter, p *Peer) {
	id := receiveThem(stream, p)

	// Test
	fmt.Println("New handshake is stored")
	them2, err := p.getPerson(id)
	if err != nil {
		panic(err)
	}

	fmt.Println("What was stored:", them2)

	receiveDrama(stream)
}

func receiveThem(stream *bufio.ReadWriter, p *Peer) string {
	// Read from stream and decode gob
	var them Person
	if err := gob.NewDecoder(stream).Decode(&them); err != nil {
		panic(err)
	}

	fmt.Println("Received handshake gob!", them)

	// Store
	p.addNewPerson(them)
	return them.ID
}

func receiveDrama(stream *bufio.ReadWriter) {

	var drama Drama
	if err := gob.NewDecoder(stream).Decode(&drama); err != nil {
		panic(err)
	}
	fmt.Println("New drama received: check it:", drama)
}

// TODO https://github.com/libp2p/specs/blob/master/discovery/mdns.md - Allows peers on same network to discover each other easily
func (p *Peer) handshake(person Person) {
	// TODO Check if this peer is already in store
	// Add person to community store
	err := p.addNewPerson(person)
	if err != nil {
		panic(err)
	}

	fmt.Println("Added " + person.ID + " to local store")

	// Should generate a shared sym key here? Prolly

	// Broadcast this to all other owned peers so they can handshake
	// If they don't ALL return a "success" msg then it failed
}
