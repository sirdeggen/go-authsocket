package wire

import (
    "encoding/hex"
    "fmt"
    "github.com/bsv-blockchain/go-sdk/bsvutil"
)

// Simple wrapper to load/serialize keys using go-sdk's bsvutil
type KeyPair struct {
    Priv *bsvutil.PrivateKey
    Pub  []byte
}

func NewKeyPairFromHex(hexpriv string) (*KeyPair, error) {
    k, err := bsvutil.NewPrivateKeyFromHex(hexpriv)
    if err != nil {
        return nil, err
    }
    pub := k.PubKey()
    return &KeyPair{Priv: k, Pub: pub}, nil
}

func (kp *KeyPair) Sign(data []byte) ([]byte, error) {
    return kp.Priv.Sign(data)
}

func (kp *KeyPair) Verify(data, sig []byte) bool {
    return kp.Pub.Verify(data, sig)
}

func (kp *KeyPair) PubKey() []byte { return kp.Pub }

func (kp *KeyPair) PubHex() string { return hex.EncodeToString(kp.Pub) }

func (kp *KeyPair) PrivHex() string { return hex.EncodeToString(kp.Priv.ToBytes()) }

func MustNewKeyPairFromHex(hexpriv string) *KeyPair { kp, _ := NewKeyPairFromHex(hexpriv); return kp }

func DemoKeypair() *KeyPair {
    // WARNING: This is a placeholder demo; replace with real key material in tests.
    // Generate 32-byte random hex string deterministically for tests is not suitable for production.
    d := "1a2b3c4d5e6f708192a3b4c5d6e7f8090a1b2c3d4e5f60718293a4b5c6d7e8f"
    k, _ := NewPrivateKeyFromHex(d)
    return &KeyPair{Priv: k, Pub: k.PubKey()}
}

// Minimal wrapper to instantiate from private key hex using go-sdk types
func NewPrivateKeyFromHex(hexpriv string) (*bsvutil.PrivateKey, error) {
    return bsvutil.NewPrivateKeyFromHex(hexpriv)
}
