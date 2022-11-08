package peer

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

const (
	ReservedWord_Owner = "me"
)

// TODO encrypt data at rest
func InitPeerStore() *bolt.DB {
	db, err := bolt.Open("locke.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Create "people" bucket and add yourself to it
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("people"))
		if err != nil {
			return fmt.Errorf("could not create 'people' bucket")
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	return db
}

// Creates a bucket for a person and initializes a relationship
// Stores peerIDs and inits a drama for each
func (p *Peer) addPerson(name string, peers []string) error {
	err := p.DB.Update(func(tx *bolt.Tx) error {
		// Create new bucket for person
		_, err := tx.CreateBucketIfNotExists([]byte(name))
		if err != nil {
			return fmt.Errorf("could not create bucket for person: %v", name)
		}

		// Init new drama for new person
		var d = CreateDrama(0)
		var drama bytes.Buffer
		err = gob.NewEncoder(&drama).Encode(d)
		if err != nil {
			return err
		}

		// Then add person to "people" bucket along with drama gob
		err = tx.Bucket([]byte("people")).Put([]byte(name), drama.Bytes())
		return err
	})

	// Add all peers to bucket and init new dramas for each
	for i := 0; i < len(peers); i++ {
		err = p.addPeer(name, peers[i])
		if err != nil {
			fmt.Print(err)
			return err
		}
	}
	return err
}

func (p *Peer) addPeer(name string, peerID string) error {
	fmt.Println("Adding peer:", peerID)

	var d = CreateDrama(0)
	var drama bytes.Buffer
	err := gob.NewEncoder(&drama).Encode(d)
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

func (p *Peer) getAllPeople() (people []Person, err error) {
	err = p.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("people"))
		b.ForEach(func(k, v []byte) error {
			person, err := p.getPerson(string(k))
			if err != nil {
				return err
			}

			people = append(people, person)
			return nil
		})

		return nil
	})
	if err != nil {
		return nil, err
	}

	return people, nil
}
