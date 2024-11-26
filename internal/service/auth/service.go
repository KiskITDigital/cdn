package auth

import (
	"time"

	"gitlab.ubrato.ru/ubrato/cdn/internal/lib/token"
)

type Service struct {
	tokenAuthorizer TokenAuthorizer
}

type TokenAuthorizer interface {
	GenerateToken(payload token.Payload) (string, error)
	GetRefreshTokenDurationLifetime() time.Duration
	ValidateToken(rawToken string) (token.Claims, error)
}

func New(
	tokenAuthorizer TokenAuthorizer,

) *Service {
	return &Service{
		tokenAuthorizer: tokenAuthorizer,
	}
}
