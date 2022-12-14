package peer

import (
	"context"
	"crypto/rand"

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

const (
	DefaultPort = "5533"
)

type Peer struct {
	Host host.Host
	DB   *bolt.DB

	// A person is defined by all the peers they own and all the relationships those peers have
	// Therefore your "self" is all your peers and their relationships
	Self
}

// type Identity struct {
// 	PeerID  string
// 	PrivKey crypto.PrivKey `json:",omitempty"`
// }

func (p *Peer) Start(name string) {
	// See if we can get a private key from store to restart a peer
	db, err := InitPeerStore()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	p.DB = db

	self, err := p.getSelf()
	if err != nil {
		fmt.Println("No existing database found, initializing a brand new peer")
		p.New(name)
		return
	}

	if self.ID != name {
		fmt.Println("Name did not match existing peer's owner")
		return
	}

	fmt.Println("Starting existing host...")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create new host
	host := createHost(ctx, self.PrivateKey, DefaultPort)
	addr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/udp/%s/quic", DefaultPort))
	if err != nil {
		panic(err)
	}

	host.Network().Listen(addr)
	defer host.Close()
	p.Host = host

	// Start listening for handshakes
	p.listenForHandshake()

	select {}
	// p.listenForCoordination()

}

// TODO I've been very liberal using *Peer as an interface for functions
// At some point a legit peer interface should be defined
// And any function that doesn't make sense to be in that interface should instead take *Peer as an argument
func (p *Peer) New(name string) {
	fmt.Println("Initializing new peer...")

	// Get flags
	// dest := flag.String("dest", "", "Destination multiaddr string")
	// peerID := flag.String("peer", "", "Peer ID")
	// port := flag.String("port", DefaultPort, "Port")
	// flag.Parse()

	// Generate key pair
	priv, _, err := crypto.GenerateECDSAKeyPair(rand.Reader)
	if err != nil {
		panic(err)
	}

	fmt.Println("Creating new host...")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create new host
	host := createHost(ctx, priv, DefaultPort)

	// New top-level drama
	tld := CreateDrama()

	// Add self to peer store
	// TODO onboarding and coordinating all peers and shizz
	peers := make(map[string]Drama)
	peers[host.ID().String()] = Drama{} // TODO this is the top-level drama so shouldnt be made new every time

	// Create self!
	p.Self = Self{
		Person: Person{
			ID:    name,
			Peers: peers,
		},
		PrivateKey: priv,
		TLD:        tld,
	}

	// Store
	err = p.addSelf(&p.Self)
	if err != nil {
		panic(err)
	}

	addr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/udp/%s/quic", DefaultPort))
	if err != nil {
		panic(err)
	}

	host.Network().Listen(addr)
	defer host.Close()
	p.Host = host

	// Start listening for handshakes
	p.listenForHandshake()
	// p.listenForCoordination()

	// Connect if dest and peerID are supplied arguments
	// if *dest != "" && *peerID != "" {
	// 	connect(ctx, &p, *dest, *peerID)
	// }

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
