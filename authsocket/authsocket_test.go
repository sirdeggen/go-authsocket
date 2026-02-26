package authsocket

import (
	"context"
	"testing"
	"time"

	"github.com/sirdeggen/go-authsocket/authsocket/transport"
	"github.com/sirdeggen/go-authsocket/internal/wire"
)

func TestAuthSocketClientServer(t *testing.T) {
	hexpriv := "0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20"
	wallet, err := wire.NewKeyPairFromHex(hexpriv)
	if err != nil {
		t.Fatal(err)
	}

	clientTransport, serverTransport := transport.InMemoryPair()

	client := NewAuthSocketClient(clientTransport, wallet)
	server := NewAuthSocketServer(serverTransport, wallet) // Server wallet not used in this test

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Start server accepting in background
	go func() {
		err := server.AcceptClient(ctx, serverTransport)
		if err != nil {
			t.Log("server accept client:", err)
		}
	}()

	// Connect client
	err = client.Connect(ctx)
	if err != nil {
		t.Fatal("client connect:", err)
	}

	// Test emit from client
	err = client.Emit(ctx, "test-event", "hello world")
	if err != nil {
		t.Fatal("client emit:", err)
	}

	// Wait a bit for message to propagate
	time.Sleep(100 * time.Millisecond)

	// Test emit from server
	err = server.Emit(ctx, "server-event", 42)
	if err != nil {
		t.Fatal("server emit:", err)
	}

	t.Log("AuthSocket client-server test completed")
}