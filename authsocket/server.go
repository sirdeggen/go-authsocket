package authsocket

import (
	"github.com/sirdeggen/go-authsocket/internal/wire"
)

// Server is a minimal stand-in for a mutual-authsocket server.
// It currently issues a nonce on Hello and accepts an Auth frame,
// returning an OK frame to signal success. This is intentionally lightweight
// to enable early integration tests against the TS client.
type Server struct{}

func NewServer() *Server { return &Server{} }

// HandleHello receives a Hello frame from a client and returns a Nonce frame.
func (s *Server) HandleHello(frame *wire.Frame) *wire.Frame {
	// In a full implementation, you would verify the client's pubkey here.
	nonce := wire.MakeNonce()
	return wire.MustFrameNonce(nonce)
}

// HandleAuth processes an Auth frame and returns an OK frame on success.
func (s *Server) HandleAuth(frame *wire.Frame) *wire.Frame {
	// In a full implementation, you would verify the signature against the nonce.
	return wire.MustFrameOK()
}
