package peer

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

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
	PeerID  string
	PrivKey string `json:",omitempty"`
}

type KeyGenerateSettings struct {
	Algorithm string
	Size      int
}

func (p *Peer) Init() {
	fmt.Println("Initializing new peer...")
	peer := Peer{}
	identity, err := CreateIdentity(KeyGenerateSettings{Algorithm: "ed25519", Size: -1})
	if err != nil {
		fmt.Println("Error!", err)
	}

	peer.Identity = identity
	fmt.Println("New peer initialized! Identity:", peer.Identity)
}

// CreateIdentity initializes a new identity.
func CreateIdentity(opts KeyGenerateSettings) (Identity, error) {
	ident := Identity{}

	var sk PrivKey
	var pk PubKey

	switch opts.Algorithm {
	case "ed25519":
		if opts.Size != -1 {
			return ident, fmt.Errorf("number of key bits does not apply when using ed25519 keys")
		}
		priv, pub, err := GenerateEd25519Key(rand.Reader)
		if err != nil {
			return ident, err
		}

		sk = priv
		pk = pub
	// case "kyber":
	// 	priv, pub, err := kyberk2so.KemKeypair768()
	// 	if err != nil {
	// 		return ident, err
	// 	}

	// 	sk = priv
	// 	pk = pub
	default:
		return ident, fmt.Errorf("unrecognized key type: %s", opts.Algorithm)
	}

	// currently storing key unencrypted. in the future we need to encrypt it.
	// TODO(security)
	skbytes, err := MarshalPrivateKey(sk)
	if err != nil {
		return ident, err
	}
	ident.PrivKey = base64.StdEncoding.EncodeToString(skbytes)

	id, err := IDFromPublicKey(pk)
	if err != nil {
		return ident, err
	}

	ident.PeerID = base58.Encode([]byte(id))
	return ident, nil
}

// IDFromPublicKey returns the Peer ID corresponding to the public key pk.
func IDFromPublicKey(pk PubKey) (string, error) {
	b, err := MarshalPublicKey(pk)
	if err != nil {
		return "", err
	}
	var alg uint64 = SHA2_256
	// if AdvancedEnableInlining && len(b) <= maxInlineKeyLength {
	// 	alg = mh.IDENTITY
	// }
	hash, _ := Sum(b, alg, -1)
	return string(hash), nil
}

func (p *Peer) Bootstrap() {
	fmt.Println("Initializing new peer")
}

func (p *Peer) Stop() {
	fmt.Println("Stopping peer: ", p.Identity.PeerID)
}
