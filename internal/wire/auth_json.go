package wire

import (
    "encoding/json"
)

type AuthMessage struct {
    Version      string      `json:"version"`
    Type         string      `json:"type"`
    IdentityKey  string      `json:"identityKey,omitempty"`
    Nonce        []int       `json:"nonce,omitempty"`
    Payload      []int       `json:"payload,omitempty"`
    Signature    string      `json:"signature,omitempty"`
    Certificates interface{} `json:"certificates,omitempty"`
}

func (a *AuthMessage) MarshalJSON() ([]byte, error) {
    type Alias AuthMessage
    return json.Marshal((*Alias)(a))
}

func (a *AuthMessage) UnmarshalJSON(data []byte) error {
    type Alias AuthMessage
    aux := &Alias{}
    if err := json.Unmarshal(data, &aux); err != nil {
        return err
    }
    *a = AuthMessage(*aux)
    return nil
}
