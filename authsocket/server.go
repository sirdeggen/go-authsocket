package authsocket

import (
    "encoding/json"
    "fmt"
    "github.com/sirdeggen/go-authsocket/internal/wire"
)

type Server struct{}

func NewServer() *Server { return &Server{} }

// HandleHello processes a Hello message (AuthMessage JSON) and responds with a nonce as number[] payload
func (s *Server) HandleHello(raw []byte) ([]byte, error) {
    var am wire.AuthMessage
    if err := json.Unmarshal(raw, &am); err != nil {
        return nil, err
    }
    if am.Type != "hello" {
        return nil, fmt.Errorf("unexpected message type: %s", am.Type)
    }
    nonce := wire.MakeNonceIntArray()
    resp := wire.AuthMessage{Version: "1", Type: "nonce", Payload: nonce}
    return json.Marshal(resp)
}

// HandleAuth processes an Auth message and returns an OK message on success
func (s *Server) HandleAuth(raw []byte) ([]byte, error) {
    var am wire.AuthMessage
    if err := json.Unmarshal(raw, &am); err != nil {
        return nil, err
    }
    if am.Type != "auth" {
        return nil, fmt.Errorf("unexpected message type: %s", am.Type)
    }
    ok := wire.AuthMessage{Version: "1", Type: "ok"}
    return json.Marshal(ok)
}
