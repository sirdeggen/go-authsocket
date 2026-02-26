package wire

import (
    "crypto/rand"
)

// MakeNonceIntArray generates a random 32-byte nonce and returns it as []int (0-255).
func MakeNonceIntArray() []int {
    b := make([]byte, 32)
    _, _ = rand.Read(b)
    out := make([]int, len(b))
    for i := range b {
        out[i] = int(b[i])
    }
    return out
}
