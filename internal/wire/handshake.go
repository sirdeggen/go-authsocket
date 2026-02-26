package wire

import (
    "crypto/rand"
    "crypto/sha256"
)

// Minimal handshake payload generator using a nonce
func MakeNonce() []byte {
    b := make([]byte, 32)
    _, _ = rand.Read(b)
    sum := sha256.Sum256(b)
    return sum[:]
}
