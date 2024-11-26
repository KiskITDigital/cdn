package contextor

import (
	"context"

	"gitlab.ubrato.ru/ubrato/cdn/internal/models"
)

func GetUserID(ctx context.Context) int {
	v, ok := ctx.Value(models.UserIDKey).(int)
	if !ok {
		return 0
	}

	return v
}

func GetOrganizationID(ctx context.Context) int {
	v, ok := ctx.Value(models.OrganizationIDKey).(int)
	if !ok {
		return 0
	}

	return v
}

func GetRole(ctx context.Context) int {
	v, ok := ctx.Value(models.RoleKey).(int)
	if !ok {
		return 0
	}

	return v
}
