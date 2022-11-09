package peer

import (
	"context"
	"crypto/rand"

	"flag"
	"fmt"
	"time"

	"github.com/boltdb/bolt"
	"github.com/multiformats/go-multiaddr"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/routing"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/p2p/net/connmgr"
	libp2pquic "github.com/libp2p/go-libp2p/p2p/transport/quic"
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
	Host host.Host
	DB   *bolt.DB

	// A person is defined by all the peers they own and all the relationships those peers have
	// Therefore you are all your peers and their relationships
	Me Person
}

// type Identity struct {
// 	PeerID  string
// 	PrivKey crypto.PrivKey `json:",omitempty"`
// }

func (p *Peer) New() {
	fmt.Println("Initializing new peer...")

	// Get flags
	dest := flag.String("dest", "", "Destination multiaddr string")
	peerID := flag.String("peer", "", "Peer ID")
	port := flag.String("port", "5533", "Port")
	name := flag.String("name", "connor", "Your name")
	flag.Parse()

	// Generate key pair
	priv, _, err := crypto.GenerateECDSAKeyPair(rand.Reader)
	if err != nil {
		panic(err)
	}

	// TODO encrypt boltdb
	p.DB = InitPeerStore()
	defer p.DB.Close()

	fmt.Println("Creating new host...")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create new host
	host := createHost(ctx, priv, *port)

	// Add self to peer store
	p.Me.ID = *name
	// p.Me.Peers[host.ID().String()] = CreateDrama(0)
	p.addNewPerson(p.Me)

	fmt.Println("Yes! I am a person:", p.Me)

	// Start listening on that host
	host.SetStreamHandler("/locke/1.0.0", p.handleHandshake)

	addr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/udp/%s/quic", *port))
	if err != nil {
		panic(err)
	}

	host.Network().Listen(addr)
	defer host.Close()
	p.Host = host

	// Connect if dest and peerID are supplied arguments
	if *dest != "" && *peerID != "" {
		connect(ctx, p, *dest, *peerID)
	}

	// Keep alive
	select {}
}

func createHost(ctx context.Context, priv crypto.PrivKey, port string) host.Host {
	connmgr, err := connmgr.NewConnManager(
		100, // Lowwater
		400, // HighWater,
		connmgr.WithGracePeriod(time.Minute),
	)
	if err != nil {
		panic(err)
	}

	host, err := libp2p.New(
		libp2p.Identity(priv),
		// support QUIC
		libp2p.Transport(libp2pquic.NewTransport(priv, nil, nil, nil)),
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
	fmt.Printf("Peer ID is: %s, address is: %s", host.ID(), host.Addrs())

	return host
}
