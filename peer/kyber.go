package peer

// import (
// 	"bytes"
// 	"crypto/subtle"

// 	kyberk2so "github.com/symbolicsoft/kyber-k2so"
// )

// // Ed25519PrivateKey is an ed25519 private key.
// type Kyber768PrivateKey struct {
// 	k [1184]byte
// }

// // Ed25519PublicKey is an ed25519 public key.
// type Kyber768PublicKey struct {
// 	k [2400]byte
// }

// // GenerateEd25519Key generates a new ed25519 private and public key pair.
// func GenerateKyber768Key() (PrivKey, PubKey, error) {
// 	pub, priv, err := kyberk2so.KemKeypair768()
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	return &Kyber768PrivateKey{
// 			k: priv,
// 		},
// 		&Kyber768PublicKey{
// 			k: pub,
// 		},
// 		nil
// }

// // Type of the private key (Kyber768).
// func (k *Kyber768PrivateKey) Type() KeyType {
// 	return KeyType_Kyber768
// }

// // Raw private key bytes.
// func (k *Kyber768PrivateKey) Raw() ([]byte, error) {
// 	buf := make([]byte, len(k.k))
// 	copy(buf, k.k[:])

// 	return buf, nil
// }

// func (k *Kyber768PrivateKey) Equals(o Key) bool {
// 	kyber, ok := o.(*Kyber768PrivateKey)
// 	if !ok {
// 		return basicEquals(k, o)
// 	}

// 	return subtle.ConstantTimeCompare(k.k[:], kyber.k[:]) == 1
// }

// // GetPublic returns an ed25519 public key from a private key.
// func (k *Kyber768PrivateKey) GetPublic() PubKey {
// 	return &Kyber768PublicKey{k: k.pubKeyBytes()}
// }

// /*
// 	Public Key
// */

// // Type of the public key (Kyber768).
// func (k *Kyber768PublicKey) Type() KeyType {
// 	return KeyType_Kyber768
// }

// // Raw public key bytes.
// func (k *Kyber768PublicKey) Raw() ([]byte, error) {
// 	return k.k[:], nil
// }

// // Equals compares two Kyber768 public keys.
// func (k *Kyber768PublicKey) Equals(o Key) bool {
// 	edk, ok := o.(*Kyber768PublicKey)
// 	if !ok {
// 		return basicEquals(k, o)
// 	}

// 	return bytes.Equal(k.k[:], edk.k[:])
// }

// // Verify checks a signature against the input data.
// // TODO
// func (k *Kyber768PublicKey) Verify(data []byte, sig []byte) (success bool, err error) {
// 	defer func() {
// 		catch.HandlePanic(recover(), &err, "ed15519 signature verification")

// 		// To be safe.
// 		if err != nil {
// 			success = false
// 		}
// 	}()
// 	return ed25519.Verify(k.k, data, sig), nil
// }

// // basicEquals
// func basicEquals(k1, k2 Key) bool {
// 	if k1.Type() != k2.Type() {
// 		return false
// 	}

// 	a, err := k1.Raw()
// 	if err != nil {
// 		return false
// 	}
// 	b, err := k2.Raw()
// 	if err != nil {
// 		return false
// 	}
// 	return subtle.ConstantTimeCompare(a, b) == 1
// }
