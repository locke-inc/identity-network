package peer

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

// TODO encrypt data at rest
func initPeerStore() *bolt.DB {
	db, err := bolt.Open("locke.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

// Creates a bucket for a person and initializes a relationship
// Stores peerIDs and inits a drama for each
func (p *Peer) addPerson(name string, peers []string) error {
	// Create bucket with person name and init new drama
	err := p.DB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(name))
		if err != nil {
			return fmt.Errorf("could not create bucket for person: %v", name)
		}

		var d = CreateDrama(0)

		// Encode new drama to gob in order to store
		var drama bytes.Buffer
		enc := gob.NewEncoder(&drama)
		err = enc.Encode(d)
		if err != nil {
			return err
		}

		err = b.Put([]byte("relationship"), drama.Bytes())
		return err
	})

	// Add all peers to bucket and init new dramas for each
	for i := 0; i < len(peers); i++ {
		err = addPeer(p, name, peers[i])
		if err != nil {
			fmt.Print(err)
			return err
		}
	}
	return err
}

func addPeer(p *Peer, name string, peerID string) error {
	fmt.Println("Adding peer:", peerID)
	var d = CreateDrama(0)

	var drama bytes.Buffer
	enc := gob.NewEncoder(&drama)
	err := enc.Encode(d)
	if err != nil {
		return err
	}

	err = p.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(name))
		err := b.Put([]byte(peerID), drama.Bytes())
		return err
	})

	return err
}

// GetPerson lists all the peerIDs owned by a person who goes by "name"
// The key is the peerID, the value is their blockchain message history (drama)
func (p *Peer) getPerson(name string) (Person, error) {
	person := Person{}
	person.Peers = make(map[string]Drama)

	p.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(name))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {

			// Decode drama from gob
			buf := bytes.NewBuffer(v)
			dec := gob.NewDecoder(buf)

			var drama Drama
			if err := dec.Decode(&drama); err != nil {
				return err
			}

			fmt.Println("Got gob!!", drama)
			person.Peers[string(k)] = drama
		}

		return nil
	})

	return person, nil
}

func blockchainMessage() {
	// Peers handshake
	// How do they identify the person? "Owner" = LAN dht

	// Two devices are advertising at the same time and "pair"
	// They send each other their owner DHT
	// Person -> []PeerID -> Blockchain message
	// Person is populated and a startRelationship() request is sent to every peerID
	// newRelationshipDrama() a drama being the "receipt" or the "script" of how the relationship plays out
	// So every peer in your owner and their owner DHT have a drama with each other
	//

	// Send identifyYourself(zk proof)
}
