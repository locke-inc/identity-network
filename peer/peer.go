package peer

import (
	"encoding/base64"
	"fmt"

	"github.com/cloudflare/circl/kem/hybrid"
	"github.com/cloudflare/circl/sign/eddilithium3"
	"github.com/ipfs/go-namesys"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/mr-tron/base58/base58"
)

/*
 1. Each peer bootstraps itself with a Gateway node
 2. Each peer has a DHT list of people it knows
 3. Each peer has a NameSystem to enable human readable usernames
 4. Each peer can handshake with other peers using a private key made by
    Crystals Kyber for e2e encryption between them. This looks the same for
    lookout nodes as it does for personal nodes.
 5. Each peer participates in community auth (peers coming to a consensus
    of how likely a given peer is to be who they claim to be). MAYBE?: Each
    auth attempt has a blockchain (essentially a merkle tree) as the message
    structure to create immutable receipts.
*/
type Peer struct {
	Identity
	DHT     *dht.IpfsDHT
	Namesys namesys.NameSystem
}

type Identity struct {
	PeerID         string
	SignPrivKey    string `json:",omitempty"`
	EncryptPrivKey string `json:",omitempty"`
}

type KeyGenerateSettings struct {
	Algorithm string
	Size      int
}

func New() Peer {
	fmt.Println("Initializing new peer...")
	peer := Peer{}
	identity, err := CreateIdentity()
	if err != nil {
		fmt.Println("Error!", err)
	}

	peer.Identity = identity
	return peer
}

// CreateIdentity initializes a new identity
// currently storing key unencrypted. in the future we need to encrypt it.
// TODO(security)
func CreateIdentity() (Identity, error) {
	ident := Identity{}

	// Generate signing keys
	signPubKey, signPrivKey, err := eddilithium3.GenerateKey(nil)
	if err != nil {
		panic(err)
	}

	packedSignPrivKey, err := signPrivKey.MarshalBinary()
	if err != nil {
		panic(err)
	}

	// Generate encryption keys
	kyber768 := hybrid.Kyber768X448()
	encryptPubKey, encryptPrivKey, err := kyber768.GenerateKeyPair()
	if err != nil {
		panic(err)
	}

	packedEncryptPrivKey, err := encryptPrivKey.MarshalBinary()
	if err != nil {
		panic(err)
	}

	ident.SignPrivKey = base64.StdEncoding.EncodeToString(packedSignPrivKey)
	ident.EncryptPrivKey = base64.StdEncoding.EncodeToString(packedEncryptPrivKey)

	// Pack public keys and create peerID
	packedSignPubKey, err := signPubKey.MarshalBinary()
	if err != nil {
		panic(err)
	}

	packedEncryptPubKey, err := encryptPubKey.MarshalBinary()
	if err != nil {
		panic(err)
	}

	id, err := IDFromPublicKey(append(packedSignPubKey, packedEncryptPubKey...))
	if err != nil {
		return ident, err
	}

	ident.PeerID = base58.Encode([]byte(id))
	return ident, nil
}

// IDFromPublicKey returns the Peer ID corresponding to the public key pk.
func IDFromPublicKey(b []byte) (string, error) {
	var alg uint64 = SHA2_256
	hash, _ := Sum(b, alg, -1)
	return string(hash), nil
}

func (p *Peer) Bootstrap() {
	fmt.Println("Initializing new peer")
}

func (p *Peer) Stop() {
	fmt.Println("Stopping peer: ", p.Identity.PeerID)
}
