package authsocket

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/sirdeggen/go-authsocket/authsocket/transport"
	"github.com/sirdeggen/go-authsocket/internal/wire"
)

func TestTransportHandshake(t *testing.T) {
	hexpriv := "0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20"
	wallet, err := wire.NewKeyPairFromHex(hexpriv)
	if err != nil {
		t.Fatal(err)
	}

	clientT, serverT := transport.InMemoryPair()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	var serverErr, clientErr error

	// Run server in background
	wg.Add(1)
	go func() {
		defer wg.Done()
		serverErr = RunServerHandshake(ctx, serverT)
	}()

	// Run client in foreground
	wg.Add(1)
	go func() {
		defer wg.Done()
		clientErr = RunClientHandshake(ctx, clientT, wallet)
	}()

	wg.Wait()

	if serverErr != nil {
		t.Fatalf("server error: %v", serverErr)
	}
	if clientErr != nil {
		t.Fatalf("client error: %v", clientErr)
	}

	t.Log("transport-based handshake completed successfully: hello -> nonce -> auth -> ok")
}

func TestTransportHandshakeTimeout(t *testing.T) {
	hexpriv := "0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20"
	wallet, err := wire.NewKeyPairFromHex(hexpriv)
	if err != nil {
		t.Fatal(err)
	}

	clientT, _ := transport.InMemoryPair()
	// Very short timeout — no server running, should fail
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err = RunClientHandshake(ctx, clientT, wallet)
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
	t.Logf("got expected error: %v", err)
}

func TestTransportMultipleHandshakes(t *testing.T) {
	for i := 0; i < 5; i++ {
		hexpriv := "0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20"
		wallet, err := wire.NewKeyPairFromHex(hexpriv)
		if err != nil {
			t.Fatal(err)
		}

		clientT, serverT := transport.InMemoryPair()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		var wg sync.WaitGroup
		var serverErr, clientErr error

		wg.Add(2)
		go func() {
			defer wg.Done()
			serverErr = RunServerHandshake(ctx, serverT)
		}()
		go func() {
			defer wg.Done()
			clientErr = RunClientHandshake(ctx, clientT, wallet)
		}()

		wg.Wait()
		cancel()

		if serverErr != nil {
			t.Fatalf("iteration %d: server error: %v", i, serverErr)
		}
		if clientErr != nil {
			t.Fatalf("iteration %d: client error: %v", i, clientErr)
		}
	}
	t.Log("5 consecutive handshakes completed successfully")
}
