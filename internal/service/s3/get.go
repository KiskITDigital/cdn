package s3

import (
	"context"

	"github.com/minio/minio-go/v7"
	"gitlab.ubrato.ru/ubrato/cdn/internal/lib/cerr"
	"gitlab.ubrato.ru/ubrato/cdn/internal/models"
)

func (s *Service) GetFile(ctx context.Context, objectName string) (*minio.Object, minio.ObjectInfo, error) {
	info, err := s.minio.StatObject(ctx, models.BucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return nil, minio.ObjectInfo{}, cerr.Wrap(err, cerr.CodeNotFound, "file not found", nil)
		}

		return nil, minio.ObjectInfo{}, err
	}

	object, err := s.minio.GetObject(context.Background(), models.BucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return nil, minio.ObjectInfo{}, cerr.Wrap(err, cerr.CodeNotFound, "file not found", nil)
		}

		return nil, minio.ObjectInfo{}, err
	}

	return object, info, nil
}
