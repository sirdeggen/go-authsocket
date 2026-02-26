package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirdeggen/go-authsocket/authsocket"
	"github.com/sirdeggen/go-authsocket/authsocket/transport"
	"github.com/sirdeggen/go-authsocket/internal/wire"
)

func main() {
	// Server wallet (dummy, as server doesn't sign in this test)
	wallet, _ := wire.NewKeyPairFromHex("030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f22")

	server := authsocket.NewAuthSocketServer(nil, wallet) // Transport set per connection

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("upgrade error:", err)
			return
		}
		defer conn.Close()

		wsTransport := &webSocketTransport{conn: conn}

		ctx := context.Background()
		err = server.AcceptClient(ctx, wsTransport)
		if err != nil {
			log.Println("accept client error:", err)
			return
		}

		// Broadcast a welcome message
		server.Emit(ctx, "message", map[string]string{"from": "Server", "text": "Welcome!"})

		// Listen for messages (in a real app, handle properly)
		for {
			data, err := wsTransport.Receive(ctx)
			if err != nil {
				break
			}
			// Parse and broadcast
			// For simplicity, just log
			log.Println("Received message:", string(data))
		}
	})

	fmt.Println("Starting authsocket server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// webSocketTransport implements transport.Transport
type webSocketTransport struct {
	conn *websocket.Conn
}

func (w *webSocketTransport) Send(ctx context.Context, data []byte) error {
	return w.conn.WriteMessage(websocket.TextMessage, data)
}

func (w *webSocketTransport) Receive(ctx context.Context) ([]byte, error) {
	_, message, err := w.conn.ReadMessage()
	return message, err
}