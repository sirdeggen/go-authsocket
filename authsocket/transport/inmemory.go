package transport

import (
	"context"
)

type Transport interface {
	Send(ctx context.Context, data []byte) error
	Receive(ctx context.Context) ([]byte, error)
}

type inMemoryClient struct {
	toServer   chan []byte
	fromServer chan []byte
}

type inMemoryServer struct {
	toClient   chan []byte
	fromClient chan []byte
}

func InMemoryPair() (Transport, Transport) {
	c2s := make(chan []byte, 1)
	s2c := make(chan []byte, 1)

	client := &inMemoryClient{toServer: c2s, fromServer: s2c}
	server := &inMemoryServer{toClient: s2c, fromClient: c2s}

	return client, server
}

func (c *inMemoryClient) Send(ctx context.Context, data []byte) error {
	select {
	case c.toServer <- data:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *inMemoryClient) Receive(ctx context.Context) ([]byte, error) {
	select {
	case v := <-c.fromServer:
		return v, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (s *inMemoryServer) Send(ctx context.Context, data []byte) error {
	select {
	case s.toClient <- data:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (s *inMemoryServer) Receive(ctx context.Context) ([]byte, error) {
	select {
	case v := <-s.fromClient:
		return v, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
