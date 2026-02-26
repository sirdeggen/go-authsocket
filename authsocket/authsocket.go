package authsocket

import (
	"context"
	"encoding/json"
	"errors"
	"sync"

	"github.com/sirdeggen/go-authsocket/authsocket/transport"
	"github.com/sirdeggen/go-authsocket/internal/wire"
)

// AuthSocketClient mimics the TypeScript AuthSocket client.
// It wraps a transport, performs handshake, and handles events.
type AuthSocketClient struct {
	transport    transport.Transport
	wallet       *wire.KeyPair
	handshaked   bool
	eventMutex   sync.RWMutex
	eventHandlers map[string][]func(data interface{})
}

// NewAuthSocketClient creates a new client with the given transport and wallet.
func NewAuthSocketClient(transport transport.Transport, wallet *wire.KeyPair) *AuthSocketClient {
	return &AuthSocketClient{
		transport:     transport,
		wallet:        wallet,
		eventHandlers: make(map[string][]func(data interface{})),
	}
}

// Connect performs the handshake over the transport.
func (c *AuthSocketClient) Connect(ctx context.Context) error {
	if c.handshaked {
		return nil
	}

	err := RunClientHandshake(ctx, c.transport, c.wallet)
	if err != nil {
		return err
	}

	c.handshaked = true

	// Start listening for incoming messages
	go c.listenForMessages(ctx)

	return nil
}

// On registers an event handler.
func (c *AuthSocketClient) On(event string, handler func(data interface{})) {
	c.eventMutex.Lock()
	defer c.eventMutex.Unlock()
	c.eventHandlers[event] = append(c.eventHandlers[event], handler)
}

// Emit sends an event with data.
func (c *AuthSocketClient) Emit(ctx context.Context, event string, data interface{}) error {
	if !c.handshaked {
		return ErrNotConnected
	}

	msg := wire.AuthMessage{
		Version: "1",
		Type:    "general",
		Payload: []int{}, // Empty for now, could encode data
	}
	// For simplicity, encode event and data as JSON in payload
	payloadData, err := json.Marshal(map[string]interface{}{"event": event, "data": data})
	if err != nil {
		return err
	}
	msg.Payload = IntsFromBytes(payloadData)

	raw, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return c.transport.Send(ctx, raw)
}

func (c *AuthSocketClient) listenForMessages(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			data, err := c.transport.Receive(ctx)
			if err != nil {
				// Handle error
				continue
			}

			var msg wire.AuthMessage
			if err := json.Unmarshal(data, &msg); err != nil {
				continue
			}

			if msg.Type == "general" && len(msg.Payload) > 0 {
				payloadBytes := BytesFromIntArray(msg.Payload)
				var eventData map[string]interface{}
				if err := json.Unmarshal(payloadBytes, &eventData); err != nil {
					continue
				}

				event, ok := eventData["event"].(string)
				if !ok {
					continue
				}

				c.eventMutex.RLock()
				handlers := c.eventHandlers[event]
				c.eventMutex.RUnlock()

				for _, handler := range handlers {
					go handler(eventData["data"])
				}
			}
		}
	}
}

// AuthSocketServer mimics the TypeScript AuthSocketServer.
// It wraps a transport, performs handshake, and broadcasts events.
type AuthSocketServer struct {
	transport    transport.Transport
	wallet       *wire.KeyPair
	handshaked   bool
	clients      map[string]*clientSession
	clientsMutex sync.RWMutex
}

type clientSession struct {
	transport transport.Transport
}

func NewAuthSocketServer(transport transport.Transport, wallet *wire.KeyPair) *AuthSocketServer {
	return &AuthSocketServer{
		transport: transport,
		wallet:    wallet,
		clients:   make(map[string]*clientSession),
	}
}

// AcceptClient performs handshake with a new client and adds to clients.
func (s *AuthSocketServer) AcceptClient(ctx context.Context, clientTransport transport.Transport) error {
	err := RunServerHandshake(ctx, clientTransport)
	if err != nil {
		return err
	}

	// Add client
	s.clientsMutex.Lock()
	s.clients["client-id"] = &clientSession{transport: clientTransport} // Use actual ID
	s.clientsMutex.Unlock()

	return nil
}

// Emit broadcasts an event to all connected clients.
func (s *AuthSocketServer) Emit(ctx context.Context, event string, data interface{}) error {
	msg := wire.AuthMessage{
		Version: "1",
		Type:    "general",
		Payload: []int{},
	}
	payloadData, err := json.Marshal(map[string]interface{}{"event": event, "data": data})
	if err != nil {
		return err
	}
	msg.Payload = IntsFromBytes(payloadData)

	raw, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	s.clientsMutex.RLock()
	defer s.clientsMutex.RUnlock()

	for _, client := range s.clients {
		go client.transport.Send(ctx, raw)
	}

	return nil
}

var ErrNotConnected = errors.New("not connected")
