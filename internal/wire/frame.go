package wire

import (
    "bytes"
    "encoding/binary"
    "fmt"
)

// Frame is a simple, versioned envelope for authsocket messages.
// This is a placeholder wire format intended to be replaced once the
// exact TS protocol (BRC-103 framing) is confirmed. It is designed to
// be drop-in compatible at the Go layer once the TS format is known.
//
// Layout (big-endian):
// - Version: 2 bytes
// - Type: 2 bytes
// - PayloadLen: 4 bytes
// - Payload: variable
// - SigLen: 4 bytes
// - Signature: variable
// - PubKeyLen: 4 bytes
// - PubKey: variable
type Frame struct {
    Version   uint16
    Type      uint16
    Payload   []byte
    Signature []byte
    PubKey    []byte
}

func (f *Frame) MarshalBinary() ([]byte, error) {
    var buf bytes.Buffer
    // Write header
    if err := binary.Write(&buf, binary.BigEndian, f.Version); err != nil {
        return nil, err
    }
    if err := binary.Write(&buf, binary.BigEndian, f.Type); err != nil {
        return nil, err
    }
    // Payload
    if err := binary.Write(&buf, binary.BigEndian, uint32(len(f.Payload))); err != nil {
        return nil, err
    }
    if len(f.Payload) > 0 {
        if _, err := buf.Write(f.Payload); err != nil {
            return nil, err
        }
    }
    // Signature
    if err := binary.Write(&buf, binary.BigEndian, uint32(len(f.Signature))); err != nil {
        return nil, err
    }
    if len(f.Signature) > 0 {
        if _, err := buf.Write(f.Signature); err != nil {
            return nil, err
        }
    }
    // PubKey
    if err := binary.Write(&buf, binary.BigEndian, uint32(len(f.PubKey))); err != nil {
        return nil, err
    }
    if len(f.PubKey) > 0 {
        if _, err := buf.Write(f.PubKey); err != nil {
            return nil, err
        }
    }
    return buf.Bytes(), nil
}

func (f *Frame) UnmarshalBinary(data []byte) error {
    r := bytes.NewReader(data)
    if err := binary.Read(r, binary.BigEndian, &f.Version); err != nil {
        return err
    }
    if err := binary.Read(r, binary.BigEndian, &f.Type); err != nil {
        return err
    }
    var plen uint32
    if err := binary.Read(r, binary.BigEndian, &plen); err != nil {
        return err
    }
    if plen > 0 {
        f.Payload = make([]byte, plen)
        if _, err := r.Read(f.Payload); err != nil {
            return err
        }
    } else {
        f.Payload = nil
    }

    var slen uint32
    if err := binary.Read(r, binary.BigEndian, &slen); err != nil {
        return err
    }
    if slen > 0 {
        f.Signature = make([]byte, slen)
        if _, err := r.Read(f.Signature); err != nil {
            return err
        }
    } else {
        f.Signature = nil
    }

    var plenPub uint32
    if err := binary.Read(r, binary.BigEndian, &plenPub); err != nil {
        return err
    }
    if plenPub > 0 {
        f.PubKey = make([]byte, plenPub)
        if _, err := r.Read(f.PubKey); err != nil {
            return err
        }
    } else {
        f.PubKey = nil
    }

    // Basic sanity checks
    if f.Version == 0 {
        return fmt.Errorf("invalid frame version 0")
    }
    return nil
}

// Helpers for framing convenience
func FrameToBytes(fr *Frame) ([]byte, error) {
    b, err := fr.MarshalBinary()
    if err != nil {
        return nil, err
    }
    return b, nil
}

func BytesToFrame(b []byte) (*Frame, error) {
    fr := &Frame{}
    if err := fr.UnmarshalBinary(b); err != nil {
        return nil, err
    }
    return fr, nil
}

// Convenience: small constant frame types for MVP
const (
    FrameTypeHello   uint16 = 1
    FrameTypeNonce   uint16 = 2
    FrameTypeAuth    uint16 = 3
    FrameTypeOK      uint16 = 4
    FrameTypeError   uint16 = 0xFFFF
)

func MustFrameHello(pubkey []byte) *Frame {
    return &Frame{Version: 1, Type: FrameTypeHello, Payload: nil, PubKey: pubkey}
}

func MustFrameNonce(nonce []byte) *Frame {
    return &Frame{Version: 1, Type: FrameTypeNonce, Payload: nonce}
}

func MustFrameAuth(sig []byte, pubkey []byte) *Frame {
    return &Frame{Version: 1, Type: FrameTypeAuth, Payload: sig, PubKey: pubkey}
}

func MustFrameOK() *Frame {
    return &Frame{Version: 1, Type: FrameTypeOK}
}

// Pretty print for debugging
func (f *Frame) String() string {
    return fmt.Sprintf("Frame{Version:%d Type:%d Payload:%d Sig:%d PubKey:%d}", f.Version, f.Type, len(f.Payload), len(f.Signature), len(f.PubKey))
}
