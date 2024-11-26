package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gitlab.ubrato.ru/ubrato/cdn/internal/config"
)

var (
	errEmptySecret = errors.New("authorizer empty JWT secret not allowed")
)

type Claims struct {
	jwt.RegisteredClaims
	Payload
}

type TokenAuthorizer struct {
	cfg config.JWT
}

func NewTokenAuthorizer(settings config.JWT) (*TokenAuthorizer, error) {
	if settings.Secret == "" {
		return nil, errEmptySecret
	}

	return &TokenAuthorizer{cfg: settings}, nil
}

func (ta *TokenAuthorizer) GetRefreshTokenDurationLifetime() time.Duration {
	return ta.cfg.LifetimeRefresh
}

func (ta *TokenAuthorizer) GenerateToken(payload Payload) (string, error) {
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ta.cfg.LifetimeAccess)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Payload: payload,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString([]byte(ta.cfg.Secret))
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}

	return ss, nil
}

func (ta *TokenAuthorizer) ValidateToken(rawToken string) (Claims, error) {
	var claims Claims
	_, err := jwt.ParseWithClaims(rawToken, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(ta.cfg.Secret), nil
	})
	if err != nil {
		return Claims{}, fmt.Errorf("validate token: %w", err)
	}

	return claims, nil
}
