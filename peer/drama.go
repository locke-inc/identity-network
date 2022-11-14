package peer

import (
	"bytes"
	cryptorand "crypto/rand" // be explicit about not being math/rand
	"encoding/gob"
	"fmt"
	"time"

	"golang.org/x/crypto/chacha20poly1305"
)

// A Drama is essentially a history of transactions among peers (or people).
// It's a blockchain with data that's encrypted by shared keys
// Blockchains are already signed, read more about sign & encrypt:
// https://std.com/~dtd/sign_encrypt/sign_encrypt7.html
type Drama struct {
	GenesisBlock Block
	Chain        []Block
	Difficulty   int
}

func CreateDrama(difficulty int) Drama {
	genesisBlock := Block{
		Hash:      "0",
		Timestamp: time.Now(),
	}
	return Drama{
		genesisBlock,
		[]Block{genesisBlock},
		difficulty,
	}
}

// Block transactional data is encrypted by xChacha20
func (b *Drama) addBlock(t Transaction, symKey []byte) {
	// Encode transaction into gob
	var transaction bytes.Buffer
	err := gob.NewEncoder(&transaction).Encode(t)
	if err != nil {
		panic(err)
	}

	// Encrypt transaction bytes
	aead, err := chacha20poly1305.NewX(symKey)
	if err != nil {
		panic(err)
	}

	// Select a random nonce, and leave capacity for the ciphertext.
	nonce := make([]byte, aead.NonceSize(), aead.NonceSize()+len(transaction.Bytes())+aead.Overhead())
	if _, err := cryptorand.Read(nonce); err != nil {
		panic(err)
	}

	// Encrypt the message and append the ciphertext to the nonce.
	encryptedMsg := aead.Seal(nonce, nonce, transaction.Bytes(), nil)

	lastBlock := b.Chain[len(b.Chain)-1]
	newBlock := Block{
		Data:         encryptedMsg,
		Nonce:        nonce,
		PreviousHash: lastBlock.Hash,
		Timestamp:    time.Now(),
	}
	newBlock.mine(b.Difficulty)
	b.Chain = append(b.Chain, newBlock)
}

// Block transactional data is encrypted by xChacha20
func (b *Drama) addUnencryptedBlock(t Transaction) {
	// Encode transaction into gob
	var transaction bytes.Buffer
	err := gob.NewEncoder(&transaction).Encode(t)
	if err != nil {
		panic(err)
	}

	lastBlock := b.Chain[len(b.Chain)-1]
	newBlock := Block{
		Data:         transaction.Bytes(),
		PreviousHash: lastBlock.Hash,
		Timestamp:    time.Now(),
	}
	newBlock.mine(b.Difficulty)
	b.Chain = append(b.Chain, newBlock)
}

func (b Drama) isValid() bool {
	fmt.Println("Validating drama...")
	for i := range b.Chain[1:] {
		previousBlock := b.Chain[i]
		currentBlock := b.Chain[i+1]
		if currentBlock.Hash != currentBlock.calculateHash() || currentBlock.PreviousHash != previousBlock.Hash {
			return false
		}
	}
	return true
}
