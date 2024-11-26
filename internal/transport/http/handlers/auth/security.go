package auth

import (
	"context"

	api "gitlab.ubrato.ru/ubrato/cdn/api/gen"
	"gitlab.ubrato.ru/ubrato/cdn/internal/models"
)

func (h *Handler) HandleBearerAuth(ctx context.Context, operationName string, t api.BearerAuth) (context.Context, error) {
	claims, err := h.authSvc.ValidateAccessToken(ctx, t.GetToken())
	if err != nil {
		return ctx, err
	}

	ctx = context.WithValue(ctx, models.UserIDKey, claims.UserID)
	ctx = context.WithValue(ctx, models.OrganizationIDKey, claims.OrganizationID)
	ctx = context.WithValue(ctx, models.RoleKey, claims.Role)

	return ctx, nil
}
