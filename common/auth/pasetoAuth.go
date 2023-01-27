package auth

import (
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type PasetoAuth struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPasetoAuth(symmetricKey string) (Auth, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, ErrInvalidKey
	}

	auth := &PasetoAuth{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}

	return auth, nil
}

func (p *PasetoAuth) CreateToken(userId string, roles []string, duration time.Duration) (string, error) {
	payload := NewPayload(userId, roles, duration)
	return p.paseto.Encrypt(p.symmetricKey, payload, nil)
}

func (p *PasetoAuth) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	if err := p.paseto.Decrypt(token, p.symmetricKey, payload, nil); err != nil {
		return nil, ErrInvalidToken
	}

	if err := payload.Valid(); err != nil {
		return payload, err
	}

	return payload, nil
}
