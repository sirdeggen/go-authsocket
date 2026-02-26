package wire

import (
	"encoding/hex"

	ec "github.com/bsv-blockchain/go-sdk/primitives/ec"
	bsvhash "github.com/bsv-blockchain/go-sdk/primitives/hash"
)

// KeyPair wrapper using go-sdk for cryptographic primitives
type KeyPair struct {
	Priv *ec.PrivateKey
	Pub  *ec.PublicKey
}

func NewKeyPairFromHex(hexpriv string) (*KeyPair, error) {
	k, err := ec.PrivateKeyFromHex(hexpriv)
	if err != nil {
		return nil, err
	}
	return &KeyPair{Priv: k, Pub: k.PubKey()}, nil
}

// Sign signs the SHA-256 hash of the data using ECDSA
func (kp *KeyPair) Sign(data []byte) ([]byte, error) {
	h := bsvhash.Sha256(data)
	sig, err := kp.Priv.Sign(h)
	if err != nil {
		return nil, err
	}
	// Use Serialize() which produces DER format suitable for transport
	return sig.Serialize(), nil
}

// Verify verifies the ECDSA signature against the data.
// Note: PublicKey.Verify() internally hashes the data with SHA-256,
// while PrivateKey.Sign() expects a pre-computed hash.
// So Sign hashes first, and Verify takes raw data.
func (kp *KeyPair) Verify(data, sigBytes []byte) bool {
	sig, err := ec.ParseSignature(sigBytes)
	if err != nil {
		return false
	}
	return kp.Pub.Verify(data, sig)
}

func (kp *KeyPair) PubKey() []byte {
	return kp.Pub.Compressed()
}

func (kp *KeyPair) PubHex() string {
	return hex.EncodeToString(kp.Pub.Compressed())
}

func (kp *KeyPair) PrivHex() string {
	return hex.EncodeToString(kp.Priv.Serialize())
}

func MustNewKeyPairFromHex(hexpriv string) *KeyPair {
	kp, _ := NewKeyPairFromHex(hexpriv)
	return kp
}

func DemoKeypair() *KeyPair {
	// A valid 32-byte hex string
	d := "1a2b3c4d5e6f708192a3b4c5d6e7f8090a1b2c3d4e5f60718293a4b5c6d7e8f"
	kp, _ := NewKeyPairFromHex(d)
	return kp
}
