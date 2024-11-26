package auth

import (
	"context"
	"log/slog"

	"gitlab.ubrato.ru/ubrato/cdn/internal/lib/token"
)

type Handler struct {
	logger  *slog.Logger
	authSvc AuthService
}

type AuthService interface {
	ValidateAccessToken(ctx context.Context, accessToken string) (token.Claims, error)
}

func New(
	logger *slog.Logger,
	authSvc AuthService,
) *Handler {
	return &Handler{
		logger:  logger,
		authSvc: authSvc,
	}
}
