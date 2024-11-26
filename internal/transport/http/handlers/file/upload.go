package file

import (
	"context"
	"fmt"

	api "gitlab.ubrato.ru/ubrato/cdn/api/gen"
	"gitlab.ubrato.ru/ubrato/cdn/internal/models"
)

func (h *Handler) UploadPost(ctx context.Context, req *api.UploadPostReq, params api.UploadPostParams) (api.UploadPostRes, error) {
	info, err := h.s3Svc.UploadFile(
		ctx,
		models.File{
			Name: req.File.Value.Name,
			Data: req.File.Value.File,
		},
		params.IsPrivate.Value,
	)
	if err != nil {
		return nil, fmt.Errorf("upload file: %w", err)
	}

	return &api.UploadPostCreated{
		Data: api.UploadPostCreatedData{
			Key: info.Key,
		},
	}, nil
}
