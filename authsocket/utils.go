package authsocket

import "errors"

func BytesFromIntArray(a []int) []byte {
    b := make([]byte, len(a))
    for i, v := range a { b[i] = byte(v) }
    return b
}

func IntsFromBytes(b []byte) []int {
    out := make([]int, len(b))
    for i, v := range b { out[i] = int(v) }
    return out
}

var ErrInvalidHandshake = errors.New("invalid handshake")
