package file

import (
	"context"
	"log/slog"

	"github.com/minio/minio-go/v7"
	"gitlab.ubrato.ru/ubrato/cdn/internal/models"
)

type Handler struct {
	logger *slog.Logger
	s3Svc  S3Service
}

type S3Service interface {
	GetFile(ctx context.Context, objectName string) (*minio.Object, minio.ObjectInfo, error)
	UploadFile(ctx context.Context, file models.File, isPrivate bool) (minio.UploadInfo, error)
}

func New(
	logger *slog.Logger,
	s3Svc S3Service,
) *Handler {
	return &Handler{
		logger: logger,
		s3Svc:  s3Svc,
	}
}
