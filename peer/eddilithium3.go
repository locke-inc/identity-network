package peer

import (
	"crypto"

	// "github.com/cloudflare/circl/sign/eddilithium3"
	"github.com/cloudflare/circl/sign/dilithium/mode3"
	"github.com/cloudflare/circl/sign/ed448"
	libp2pcrypto "github.com/libp2p/go-libp2p/core/crypto"
	crypto_pb "github.com/libp2p/go-libp2p/core/crypto/pb"
	eddilithium3 "github.com/locke-inc/identity-network/peer/crypto"
)

type Eddilithium3PrivKey struct {
	Priv *eddilithium3.PrivateKey
}

type Eddilithium3PubKey struct {
	Pub *eddilithium3.PublicKey
}

// Private key functions to satisfy interface
// (1) Equals
func (privKey *Eddilithium3PrivKey) Equals(o libp2pcrypto.Key) bool {
	return privKey.Priv.Equal(o)
}

// (2) Raw
func (privKey *Eddilithium3PrivKey) Raw() ([]byte, error) {
	keyBytes, err := privKey.Priv.MarshalBinary()
	if err != nil {
		panic(err)
	}

	return keyBytes, nil
}

// (3) Type
func (privKey *Eddilithium3PrivKey) Type() crypto_pb.KeyType {
	return 4 // TODO add to protobuf
}

// (4) Sign
func (privKey *Eddilithium3PrivKey) Sign(data []byte) ([]byte, error) {
	return privKey.Priv.Sign(nil, data, crypto.Hash(0))
}

// (5) GetPublic
func (privKey *Eddilithium3PrivKey) GetPublic() libp2pcrypto.PubKey {
	return &Eddilithium3PubKey{
		Pub: &eddilithium3.PublicKey{
			privKey.Priv.E.Public().(ed448.PublicKey),
			*privKey.Priv.D.Public().(*mode3.PublicKey),
		},
	}
}

// Public key functions to satisfy interface
// (1) Equals
func (pubKey *Eddilithium3PubKey) Equals(o libp2pcrypto.Key) bool {
	return pubKey.Pub.Equal(o)
}

// (2) Raw
func (pubKey *Eddilithium3PubKey) Raw() ([]byte, error) {
	keyBytes, err := pubKey.Pub.MarshalBinary()
	if err != nil {
		panic(err)
	}

	return keyBytes, nil
}

// (3) Type
func (pubKey *Eddilithium3PubKey) Type() crypto_pb.KeyType {
	return 4 // TODO add to protobuf
}

// (4) Verify
func (pubKey *Eddilithium3PubKey) Verify(data, sigBytes []byte) (success bool, err error) {
	return eddilithium3.Verify(pubKey.Pub, data, sigBytes), nil
}
