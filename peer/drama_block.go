package peer

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Block struct {
	Data         []byte // Gob encoded Transaction, encrypted into bytes
	Nonce        []byte
	Hash         string
	PreviousHash string
	Timestamp    time.Time
	Pow          int
}

type Transaction struct {
	Requester   string // peerID that initiated the transaction
	RequestType string // "auth" for example

	Responder   string // peerID that responded
	Result      int    // for now let's say it's between 0 and 100%
	Application string // context for what the request is for
	ProcessTime time.Duration
}

func (b Block) calculateHash() string {
	blockData := b.PreviousHash + string(b.Data) + b.Timestamp.String() + strconv.Itoa(b.Pow)
	blockHash := sha256.Sum256([]byte(blockData))
	return fmt.Sprintf("%x", blockHash)
}

func (b *Block) mine(difficulty int) {
	for !strings.HasPrefix(b.Hash, strings.Repeat("0", difficulty)) {
		b.Pow++
		b.Hash = b.calculateHash()
	}
}
