package file

import (
	"context"
	"fmt"
	"mime"
	"path/filepath"
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

	isPrivate, err := strconv.ParseBool(info.UserMetadata[models.MetaPrivateKey])
	if err != nil {
		return nil, fmt.Errorf("convert private tag: %w", err)
	}

	ext := filepath.Ext(info.Key)

	if !isPrivate {
		return &api.FileIDGetOKHeaders{
			ContentLength: int(info.Size),
			LastModified:  info.LastModified,
			XFileType:     mime.TypeByExtension(ext),
			Response: api.FileIDGetOK{
				Data: object,
			},
		}, nil
	}

	userIDString := info.UserMetadata[models.MetaUserID]

	userID, err := strconv.Atoi(userIDString)
	if err != nil {
		return nil, fmt.Errorf("convert user id: %w", err)
	}

	requestUserID := contextor.GetUserID(ctx)
	requestUserRole := contextor.GetRole(ctx)

	if userID != requestUserID && requestUserRole < models.UserRoleEmployee {
		return nil, cerr.Wrap(cerr.ErrPermission, cerr.CodeNotPermitted, "you don't have authorization to retrieve this file", nil)
	}

	return &api.FileIDGetOKHeaders{
		ContentLength: int(info.Size),
		LastModified:  info.LastModified,
		XFileType:     mime.TypeByExtension(ext),
		Response: api.FileIDGetOK{
			Data: object,
		},
	}, nil
}

func (h *Handler) FileIDHead(ctx context.Context, params api.FileIDHeadParams) (api.FileIDHeadRes, error) {
	_, info, err := h.s3Svc.GetFile(ctx, params.ID)
	if err != nil {
		return nil, err
	}

	ext := filepath.Ext(info.Key)

	return &api.FileIDHeadOK{
		ContentLength: int(info.Size),
		LastModified:  info.LastModified,
		XFileType:     mime.TypeByExtension(ext),
	}, nil
}
