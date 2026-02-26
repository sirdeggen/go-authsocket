package transport

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// WebSocketTransport implements Transport over WebSocket for the client side.
type WebSocketTransport struct {
	conn *websocket.Conn
	mu   sync.Mutex
}

// NewWebSocketClient connects to a WebSocket server and returns a Transport.
func NewWebSocketClient(url string) (Transport, error) {
	dialer := websocket.DefaultDialer
	dialer.HandshakeTimeout = 5 * time.Second

	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		return nil, fmt.Errorf("dial websocket: %w", err)
	}

	return &WebSocketTransport{conn: conn}, nil
}

func (w *WebSocketTransport) Send(ctx context.Context, data []byte) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Set write deadline
	deadline, ok := ctx.Deadline()
	if !ok {
		deadline = time.Now().Add(5 * time.Second)
	}
	w.conn.SetWriteDeadline(deadline)

	return w.conn.WriteMessage(websocket.TextMessage, data)
}

func (w *WebSocketTransport) Receive(ctx context.Context) ([]byte, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Set read deadline
	deadline, ok := ctx.Deadline()
	if !ok {
		deadline = time.Now().Add(5 * time.Second)
	}
	w.conn.SetReadDeadline(deadline)

	_, message, err := w.conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	return message, nil
}

func (w *WebSocketTransport) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.conn.Close()
}

// Note: For full WebSocket server integration, implement a server that accepts connections
// and runs RunServerHandshake. For now, this is client-only transport.
