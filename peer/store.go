package peer

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"

	"github.com/libp2p/go-libp2p/core/crypto"

	"github.com/boltdb/bolt"
)

const (
	ReservedWord_Self = "self"
	ReservedWord_TLD  = "tld"
	Prefix_Peer       = "peer_"
	Prefix_Key        = "key_"
)

// TODO encrypt data at rest
func InitPeerStore() (*bolt.DB, error) {
	db, err := bolt.Open("locke.db", 0600, nil)
	if err != nil {
		return nil, err
	}

	// Create "people" bucket
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("people"))
		if err != nil {
			return fmt.Errorf("could not create 'people' bucket")
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Creates a bucket for a person and initializes a relationship
// Stores peerIDs and inits a drama for each
func (p *Peer) addSelf(self *Self) error {
	err := p.DB.Update(func(tx *bolt.Tx) error {
		// Add "self" to people bucket
		err := tx.Bucket([]byte("people")).Put([]byte(ReservedWord_Self), []byte(self.ID))
		if err != nil {
			return err
		}

		// Create new bucket for SELF
		b, err := tx.CreateBucketIfNotExists([]byte(ReservedWord_Self))
		if err != nil {
			return fmt.Errorf("could not create SELF bucket")
		}

		// Add private key to self bucket
		keyBytes, err := self.PrivateKey.Raw()
		if err != nil {
			return err
		}
		b.Put([]byte(Prefix_Key+"privKey"), keyBytes)

		// Add top-level drama to bucket
		var drama bytes.Buffer
		err = gob.NewEncoder(&drama).Encode(self.TLD)
		if err != nil {
			return err
		}
		b.Put([]byte(ReservedWord_TLD), drama.Bytes())

		return err
	})

	// Add all peers to person's bucket and init new dramas for each
	for peerID, d := range self.Person.Peers {
		err = p.addPeer(ReservedWord_Self, peerID, &d)
		if err != nil {
			fmt.Print(err)
			return err
		}
	}

	return err
}

// Creates a bucket for a person and initializes a relationship
// Stores peerIDs and inits a drama for each
func (p *Peer) addPerson(person *Person, d *Drama, symKey *[]byte) error {
	err := p.DB.Update(func(tx *bolt.Tx) error {
		// Create new bucket for person
		_, err := tx.CreateBucketIfNotExists([]byte(person.ID))
		if err != nil {
			return fmt.Errorf("could not create bucket for person: %v", person.ID)
		}

		var drama bytes.Buffer
		err = gob.NewEncoder(&drama).Encode(d)
		if err != nil {
			return err
		}

		// Then add person to "people" bucket along with drama gob
		err = tx.Bucket([]byte("people")).Put([]byte(person.ID), drama.Bytes())
		return err
	})

	// Add all peers to person's bucket and init new dramas for each
	for peerID, _ := range person.Peers {
		d := CreateDrama()
		err = p.addPeer(person.ID, peerID, &d)
		if err != nil {
			fmt.Print(err)
			return err
		}
	}

	// Store symKey
	// TODO figure out clever key naming system
	p.addKey(person.ID, "privKey", *symKey)
	return err
}

func (p *Peer) addPeer(name string, pid string, d *Drama) error {
	fmt.Println("Adding peer:", pid)

	var drama bytes.Buffer
	err := gob.NewEncoder(&drama).Encode(d)
	if err != nil {
		return err
	}

	err = p.DB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(name))
		if err != nil {
			return err
		}

		err = b.Put([]byte(Prefix_Peer+pid), drama.Bytes())
		return err
	})

	return err
}

func (p *Peer) addKey(personName string, keyName string, key []byte) error {
	fmt.Println("Adding key:", keyName)

	err := p.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(personName))
		err := b.Put([]byte(Prefix_Key+keyName), key)
		return err
	})

	return err
}

func (p *Peer) getSelf() (Self, error) {
	self := Self{}

	// Get Self's personame
	err := p.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("people"))
		if b == nil {
			return errors.New("People database not found")
		}

		self.ID = string(b.Get([]byte(ReservedWord_Self)))
		return nil
	})
	if err != nil {
		return Self{}, err
	}

	// Get private key, top-level drama, and peers from self bucket
	err = p.DB.View(func(tx *bolt.Tx) error {
		// Get "self" bucket
		b := tx.Bucket([]byte(ReservedWord_Self))
		if b == nil {
			return errors.New(ReservedWord_Self + " -- database not found")
		}

		// Get private key
		key := b.Get([]byte(Prefix_Key + "privKey"))
		priv, err := crypto.UnmarshalECDSAPrivateKey(key)
		if err != nil {
			return errors.New("Couldn't load private key")
		}
		self.PrivateKey = priv

		// Get top-level drama
		tld := b.Get([]byte("tld"))
		buf := bytes.NewBuffer(tld)
		dec := gob.NewDecoder(buf)

		var drama Drama
		if err := dec.Decode(&drama); err != nil {
			return err
		}

		self.TLD = drama

		// Get all peers
		self.Peers = make(map[string]Drama)
		c := b.Cursor()
		for k, v := c.Seek([]byte(Prefix_Peer)); k != nil && bytes.HasPrefix(k, []byte(Prefix_Peer)); k, v = c.Next() {
			// Decode drama from gob
			buf := bytes.NewBuffer(v)
			dec := gob.NewDecoder(buf)

			var drama Drama
			if err := dec.Decode(&drama); err != nil {
				return err
			}

			self.Person.Peers[string(k)] = drama
		}

		return nil
	})
	if err != nil {
		return Self{}, err
	}

	return self, nil
}

// GetPerson lists all the peerIDs owned by a person who goes by "name"
// The key is the peerID, the value is their blockchain message history (drama)
func (p *Peer) getPerson(name string) (Person, error) {
	person := Person{}
	person.ID = name
	person.Peers = make(map[string]Drama)

	p.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(name))
		c := b.Cursor()
		// Get all the person's peers
		for k, v := c.Seek([]byte(Prefix_Peer)); k != nil && bytes.HasPrefix(k, []byte(Prefix_Peer)); k, v = c.Next() {
			// Decode drama from gob
			buf := bytes.NewBuffer(v)
			dec := gob.NewDecoder(buf)

			var drama Drama
			if err := dec.Decode(&drama); err != nil {
				return err
			}

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

func (p *Peer) updateDrama(name string, pid string, d *Drama) error {
	fmt.Println("Updating Drama for peer:", pid)

	// TODO
	// if !d.isValid() {
	// 	return errors.New("Cannot update Drama, it is invalid")
	// }

	var drama bytes.Buffer
	err := gob.NewEncoder(&drama).Encode(d)
	if err != nil {
		return err
	}

	err = p.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(name))
		err = b.Put([]byte(Prefix_Peer+pid), drama.Bytes())
		return err
	})

	return err
}
