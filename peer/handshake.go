package peer

import (
	"github.com/cloudflare/circl/kem"
)

func handshake(pubKey kem.PublicKey) {
	// Generate shared secret with KEM
	// Store secret and use for AES
	// Send to peer

	// ct, ss, err := hybrid.Kyber768X448().Encapsulate(pubKey)
	// if err != nil {
	// 	panic(err)
	// }

	// Store ss in AddrBook
}
