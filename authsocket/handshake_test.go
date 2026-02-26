package authsocket

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	w "github.com/sirdeggen/go-authsocket/internal/wire"
)

func TestInProcessHandshake(t *testing.T) {
	s := NewServer()
	hexpriv := "0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20"
	clientKey, err := w.NewKeyPairFromHex(hexpriv)
	if err != nil {
		t.Fatal(err)
	}
	c := NewClient(clientKey)

	// 1. Client sends Hello
	hello, err := c.Hello()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("hello: %s", string(hello))

	// Server responds with nonce
	nonceRaw, err := s.HandleHello(hello)
	if err != nil {
		t.Fatal(err)
	}

	var nonceMsg w.AuthMessage
	if err := json.Unmarshal(nonceRaw, &nonceMsg); err != nil {
		t.Fatalf("nonce decode: %v", err)
	}
	if nonceMsg.Type != "nonce" {
		t.Fatalf("expected type=nonce, got %s", nonceMsg.Type)
	}
	if len(nonceMsg.Payload) != 32 {
		t.Fatalf("expected 32-byte nonce payload, got %d", len(nonceMsg.Payload))
	}
	for _, v := range nonceMsg.Payload {
		if v < 0 || v > 255 {
			t.Fatalf("nonce payload value out of byte range: %d", v)
		}
	}

	t.Logf("nonce: %v", nonceMsg.Payload)

	// 2. Client sends Auth with signed nonce
	auth, err := c.Auth(nonceMsg.Payload)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("auth: %s", string(auth))

	// Verify the auth message has correct shape
	var authMsg w.AuthMessage
	if err := json.Unmarshal(auth, &authMsg); err != nil {
		t.Fatalf("auth decode: %v", err)
	}
	if authMsg.Type != "auth" {
		t.Fatalf("expected type=auth, got %s", authMsg.Type)
	}
	if authMsg.IdentityKey == "" {
		t.Fatal("auth message missing identityKey")
	}
	if authMsg.Signature == "" {
		t.Fatal("auth message missing signature")
	}

	// Verify signature independently
	nonceBytes := make([]byte, len(nonceMsg.Payload))
	for i, v := range nonceMsg.Payload {
		nonceBytes[i] = byte(v)
	}
	sigBytes, err := hex.DecodeString(authMsg.Signature)
	if err != nil {
		t.Fatalf("failed to decode signature hex: %v", err)
	}
	if !clientKey.Verify(nonceBytes, sigBytes) {
		t.Fatal("signature verification failed")
	}

	t.Log("signature verified successfully")

	// Server responds with OK
	okRaw, err := s.HandleAuth(auth)
	if err != nil {
		t.Fatal(err)
	}

	var okMsg w.AuthMessage
	if err := json.Unmarshal(okRaw, &okMsg); err != nil {
		t.Fatalf("ok decode: %v", err)
	}
	if okMsg.Type != "ok" {
		t.Fatalf("expected type=ok, got %s", okMsg.Type)
	}

	t.Log("handshake complete: hello -> nonce -> auth -> ok")
}

func TestAuthMessageJSONRoundTrip(t *testing.T) {
	msg := w.AuthMessage{
		Version:     "1",
		Type:        "general",
		IdentityKey: "abc123",
		Payload:     []int{72, 101, 108, 108, 111},
		Signature:   "deadbeef",
	}

	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("json: %s", string(data))

	var decoded w.AuthMessage
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}

	if decoded.Version != msg.Version {
		t.Fatalf("version mismatch: %s != %s", decoded.Version, msg.Version)
	}
	if decoded.Type != msg.Type {
		t.Fatalf("type mismatch: %s != %s", decoded.Type, msg.Type)
	}
	if decoded.IdentityKey != msg.IdentityKey {
		t.Fatalf("identityKey mismatch")
	}
	if len(decoded.Payload) != len(msg.Payload) {
		t.Fatalf("payload length mismatch")
	}
	for i := range msg.Payload {
		if decoded.Payload[i] != msg.Payload[i] {
			t.Fatalf("payload[%d] mismatch: %d != %d", i, decoded.Payload[i], msg.Payload[i])
		}
	}
	if decoded.Signature != msg.Signature {
		t.Fatalf("signature mismatch")
	}
}

func TestNonceIntArray(t *testing.T) {
	nonce := w.MakeNonceIntArray()
	if len(nonce) != 32 {
		t.Fatalf("expected 32 elements, got %d", len(nonce))
	}
	for _, v := range nonce {
		if v < 0 || v > 255 {
			t.Fatalf("value out of byte range: %d", v)
		}
	}
	// Ensure two consecutive nonces differ (probabilistic but near-certain)
	nonce2 := w.MakeNonceIntArray()
	same := true
	for i := range nonce {
		if nonce[i] != nonce2[i] {
			same = false
			break
		}
	}
	if same {
		t.Fatal("two consecutive nonces are identical — RNG issue")
	}
}
