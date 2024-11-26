package file

import (
	"context"
	"fmt"
	"strconv"

	api "gitlab.ubrato.ru/ubrato/cdn/api/gen"
	"gitlab.ubrato.ru/ubrato/cdn/internal/lib/cerr"
	"gitlab.ubrato.ru/ubrato/cdn/internal/lib/contextor"
	"gitlab.ubrato.ru/ubrato/cdn/internal/models"
)

func (h *Handler) FileIDGet(ctx context.Context, params api.FileIDGetParams) (api.FileIDGetRes, error) {
	object, info, err := h.s3Svc.GetFile(ctx, params.ID)
	if err != nil {
		return nil, fmt.Errorf("get file: %w", err)
	}

	if object == nil {
		return nil, cerr.Wrap(cerr.ErrNotFound, cerr.CodeNotFound, "file not found", nil)
	}

	defer object.Close()

	fmt.Printf("%#v\n", info.UserMetadata)

	isPrivate, err := strconv.ParseBool(info.UserMetadata["Private"])
	if err != nil {
		return nil, fmt.Errorf("convert private tag: %w", err)
	}

	if !isPrivate {
		return &api.FileIDGetOK{
			Data: object,
		}, nil
	}

	userIDString := info.UserMetadata["User_id"]

	userID, err := strconv.Atoi(userIDString)
	if err != nil {
		return nil, fmt.Errorf("convert user id: %w", err)
	}

	requestUserID := contextor.GetUserID(ctx)
	requestUserRole := contextor.GetRole(ctx)

	if userID != requestUserID && requestUserRole < models.UserRoleEmployee {
		return nil, cerr.Wrap(cerr.ErrPermission, cerr.CodeNotPermitted, "you don't have authorization to retrieve this file", nil)
	}

	return &api.FileIDGetOK{
		Data: object,
	}, nil
}
