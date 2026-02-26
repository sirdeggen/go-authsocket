package authsocket

import (
	"encoding/hex"
	"encoding/json"

	"github.com/sirdeggen/go-authsocket/internal/wire"
)

// Client scaffolding for in-process handshake with above server
type Client struct{ Wallet *wire.KeyPair }

func NewClient(w *wire.KeyPair) *Client { return &Client{Wallet: w} }

func (c *Client) Hello() ([]byte, error) {
	pubHex := hex.EncodeToString(c.Wallet.PubKey())
	am := wire.AuthMessage{Version: "1", Type: "hello", IdentityKey: pubHex}
	return json.Marshal(am)
}

func (c *Client) Auth(nonce []int) ([]byte, error) {
	nonceBytes := make([]byte, len(nonce))
	for i, v := range nonce {
		nonceBytes[i] = byte(v)
	}
	sig, err := c.Wallet.Sign(nonceBytes)
	if err != nil {
		return nil, err
	}
	am := wire.AuthMessage{
		Version:     "1",
		Type:        "auth",
		Payload:     nonce,
		IdentityKey: hex.EncodeToString(c.Wallet.PubKey()),
		Signature:   hex.EncodeToString(sig),
	}
	return json.Marshal(am)
}
