package authsocket

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sirdeggen/go-authsocket/authsocket/transport"
	"github.com/sirdeggen/go-authsocket/internal/wire"
)

// RunClientHandshake drives the client side of the handshake over a transport.
// It sends Hello, waits for Nonce, sends Auth, waits for OK.
func RunClientHandshake(ctx context.Context, t transport.Transport, wallet *wire.KeyPair) error {
	c := NewClient(wallet)

	// 1. Send Hello
	hello, err := c.Hello()
	if err != nil {
		return fmt.Errorf("hello: %w", err)
	}
	if err := t.Send(ctx, hello); err != nil {
		return fmt.Errorf("send hello: %w", err)
	}

	// 2. Receive Nonce
	nonceRaw, err := t.Receive(ctx)
	if err != nil {
		return fmt.Errorf("receive nonce: %w", err)
	}
	var nonceMsg wire.AuthMessage
	if err := json.Unmarshal(nonceRaw, &nonceMsg); err != nil {
		return fmt.Errorf("decode nonce: %w", err)
	}
	if nonceMsg.Type != "nonce" {
		return fmt.Errorf("expected type=nonce, got %s", nonceMsg.Type)
	}

	// 3. Send Auth
	auth, err := c.Auth(nonceMsg.Payload)
	if err != nil {
		return fmt.Errorf("auth: %w", err)
	}
	if err := t.Send(ctx, auth); err != nil {
		return fmt.Errorf("send auth: %w", err)
	}

	// 4. Receive OK
	okRaw, err := t.Receive(ctx)
	if err != nil {
		return fmt.Errorf("receive ok: %w", err)
	}
	var okMsg wire.AuthMessage
	if err := json.Unmarshal(okRaw, &okMsg); err != nil {
		return fmt.Errorf("decode ok: %w", err)
	}
	if okMsg.Type != "ok" {
		return fmt.Errorf("expected type=ok, got %s", okMsg.Type)
	}

	return nil
}

// RunServerHandshake drives the server side of the handshake over a transport.
// It waits for Hello, sends Nonce, waits for Auth, sends OK.
func RunServerHandshake(ctx context.Context, t transport.Transport) error {
	s := NewServer()

	// 1. Receive Hello
	helloRaw, err := t.Receive(ctx)
	if err != nil {
		return fmt.Errorf("receive hello: %w", err)
	}

	// 2. Process Hello -> Send Nonce
	nonceReply, err := s.HandleHello(helloRaw)
	if err != nil {
		return fmt.Errorf("handle hello: %w", err)
	}
	if err := t.Send(ctx, nonceReply); err != nil {
		return fmt.Errorf("send nonce: %w", err)
	}

	// 3. Receive Auth
	authRaw, err := t.Receive(ctx)
	if err != nil {
		return fmt.Errorf("receive auth: %w", err)
	}

	// 4. Process Auth -> Send OK
	okReply, err := s.HandleAuth(authRaw)
	if err != nil {
		return fmt.Errorf("handle auth: %w", err)
	}
	if err := t.Send(ctx, okReply); err != nil {
		return fmt.Errorf("send ok: %w", err)
	}

	return nil
}
