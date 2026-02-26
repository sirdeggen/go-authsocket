package wire

import (
	"encoding/hex"
	"testing"
)

func TestSignVerifyRoundTrip(t *testing.T) {
	kp, err := NewKeyPairFromHex("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20")
	if err != nil {
		t.Fatal(err)
	}

	data := []byte("hello world")
	sig, err := kp.Sign(data)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("sig hex: %s", hex.EncodeToString(sig))
	t.Logf("sig len: %d", len(sig))

	if !kp.Verify(data, sig) {
		t.Fatal("direct sign/verify round trip failed")
	}

	// Also test with hex encode/decode round trip (like the JSON transport does)
	sigHex := hex.EncodeToString(sig)
	sigBack, err := hex.DecodeString(sigHex)
	if err != nil {
		t.Fatal(err)
	}
	if !kp.Verify(data, sigBack) {
		t.Fatal("hex round trip sign/verify failed")
	}

	t.Log("sign/verify round trip passed")
}
