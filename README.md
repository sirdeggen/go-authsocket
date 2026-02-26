# go-authsocket

Go port of [bsv-blockchain/authsocket](https://github.com/bsv-blockchain/authsocket), preserving the exact TypeScript JSON wire format for drop-in compatibility.

## Overview

`go-authsocket` implements the AuthSocket protocol in Go using [bsv-blockchain/go-sdk](https://github.com/bsv-blockchain/go-sdk) for all cryptographic primitives (secp256k1 ECDSA).

The wire format is a JSON `AuthMessage` object with payloads as `number[]` (integer arrays), matching the TypeScript client/server exactly.

## Wire Format

All messages are JSON objects exchanged on an `authMessage` channel:

```json
{
  "version": "1",
  "type": "hello|nonce|auth|ok",
  "identityKey": "<compressed-pubkey-hex>",
  "payload": [0, 255, 128, ...],
  "signature": "<der-signature-hex>",
  "certificates": null
}
```

### Handshake Flow

```
Client                          Server
  |                               |
  |--- hello (identityKey) ------>|
  |                               |
  |<----- nonce (payload[]) ------|
  |                               |
  |--- auth (sig, identityKey) -->|
  |                               |
  |<---------- ok ----------------|
```

## Usage

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/sirdeggen/go-authsocket/authsocket"
    "github.com/sirdeggen/go-authsocket/authsocket/transport"
    "github.com/sirdeggen/go-authsocket/internal/wire"
)

func main() {
    wallet, _ := wire.NewKeyPairFromHex("your-private-key-hex")
    clientT, serverT := transport.InMemoryPair()
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    go authsocket.RunServerHandshake(ctx, serverT)

    err := authsocket.RunClientHandshake(ctx, clientT, wallet)
    if err != nil {
        fmt.Println("handshake failed:", err)
    } else {
        fmt.Println("handshake succeeded!")
    }
}
```

## Compatibility

Designed to be wire-compatible with:
- [bsv-blockchain/authsocket](https://github.com/bsv-blockchain/authsocket) (TypeScript server)
- [bsv-blockchain/authsocket-client](https://github.com/bsv-blockchain/authsocket-client) (TypeScript client)

## Dependencies

- [bsv-blockchain/go-sdk](https://github.com/bsv-blockchain/go-sdk) v1.2.18 — secp256k1 ECDSA, SHA-256, DER signatures

## Testing

```bash
go test ./... -v
```

## License

MIT
