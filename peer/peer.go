package peer

import (
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"log"
	"time"

	// "github.com/cloudflare/circl/sign/eddilithium3"
	eddilithium3 "github.com/locke-inc/identity-network/peer/crypto"

	"github.com/ipfs/go-namesys"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/routing"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/p2p/net/connmgr"

	// "github.com/libp2p/go-libp2p/p2p/net/connmgr"

	libp2pquic "github.com/libp2p/go-libp2p/p2p/transport/quic"

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
	PrivKey *Eddilithium3PrivKey `json:",omitempty"`
}

type KeyGenerateSettings struct {
	Algorithm string
	Size      int
}

func (p *Peer) New() {
	fmt.Println("Initializing new peer...")

	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	dest := flag.String("d", "", "Destination multiaddr string")
	peerID := flag.String("p", "", "Peer ID")
	flag.Parse()

	// identity, err := createIdentity()
	// if err != nil {
	// 	fmt.Println("Error!", err)
	// }

	// p.Identity = identity
	// Generate signing keys

	priv, _, err := crypto.GenerateECDSAKeyPair(rand.Reader)
	if err != nil {
		panic(err)
	}
	// _, privKey, err := eddilithium3.GenerateKey(nil)
	// if err != nil {
	// 	panic(err)
	// }

	// sk := Eddilithium3PrivKey{privKey}

	// fmt.Println("New peer identity created, ID:", identity.PeerID)
	fmt.Println("Creating new host...")

	// host := CreateHost(p.Identity.PrivKey, ctx)
	// host := createHost(&sk, ctx)

	fmt.Println("New host created, connecting to:", *dest)

	// Starting daemon
	err = run(priv, *dest, *peerID)
	// rw, err := startPeerAndConnect(ctx, host, *dest)

	if err != nil {
		log.Println(err)
		return
	}

	// Create a thread to read and write data.
	// go writeData(rw)
	// go readData(rw)
}

// CreateIdentity initializes a new identity
// currently storing key unencrypted. in the future we need to encrypt it.
// TODO(security)
func createIdentity() (Identity, error) {
	ident := Identity{}

	// Generate signing keys
	pubKey, privKey, err := eddilithium3.GenerateKey(nil)
	if err != nil {
		panic(err)
	}

	ident.PrivKey = &Eddilithium3PrivKey{privKey}

	// Pack public keys and create peerID
	packedPubKey, err := pubKey.MarshalBinary()
	if err != nil {
		panic(err)
	}

	id, err := IDFromPublicKey(packedPubKey)
	if err != nil {
		return ident, err
	}

	ident.PeerID = base58.Encode([]byte(id))

	return ident, nil
}

func createHost(priv crypto.PrivKey, ctx context.Context) host.Host {
	connmgr, err := connmgr.NewConnManager(
		100, // Lowwater
		400, // HighWater,
		connmgr.WithGracePeriod(time.Minute),
	)
	if err != nil {
		panic(err)
	}

	host, err := libp2p.New(
		// Use the keypair we generated
		libp2p.Identity(priv),
		// Multiple listen addresses
		libp2p.ListenAddrStrings(
			"/ip4/0.0.0.0/udp/9000/quic", // a UDP endpoint for the QUIC transport
		),
		// support QUIC
		libp2p.Transport(libp2pquic.NewTransport),
		// support any other default transports (TCP)
		// libp2p.DefaultTransports,
		// Let's prevent our peer from having too many
		// connections by attaching a connection manager.
		libp2p.ConnectionManager(connmgr),
		// Attempt to open ports using uPNP for NATed hosts.
		libp2p.NATPortMap(),
		// Let this host use the DHT to find other hosts
		libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
			idht, err := dht.New(ctx, h)
			return idht, err
		}),
		// Let this host use relays and advertise itself on relays if
		// it finds it is behind NAT. Use libp2p.Relay(options...) to
		// enable active relays and more.
		libp2p.EnableAutoRelay(),
	)
	if err != nil {
		panic(err)
	}
	defer host.Close()
	fmt.Printf("Hello World, my second hosts ID is %s\n", host.ID())

	return host
}

// IDFromPublicKey returns the Peer ID corresponding to the public key pk.
func IDFromPublicKey(b []byte) (string, error) {
	var alg uint64 = SHA2_256
	hash, _ := Sum(b, alg, -1)
	return string(hash), nil
}
