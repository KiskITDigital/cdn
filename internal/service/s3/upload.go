package s3

import (
	"context"
	"fmt"
	"strconv"

	"github.com/minio/minio-go/v7"
	"gitlab.ubrato.ru/ubrato/cdn/internal/lib/contextor"
	"gitlab.ubrato.ru/ubrato/cdn/internal/models"
)

func (s *Service) UploadFile(ctx context.Context, file models.File, isPrivate bool) (minio.UploadInfo, error) {
	snowflakeID := s.node.Generate().String()
	objectName := fmt.Sprintf("%s_%s", snowflakeID, file.Name)

	info, err := s.minio.PutObject(context.Background(), models.BucketName, objectName, file.Data, -1, minio.PutObjectOptions{
		UserMetadata: map[string]string{
			models.MetaUserID:     strconv.FormatInt(int64(contextor.GetUserID(ctx)), 10),
			models.MetaPrivateKey: strconv.FormatBool(isPrivate),
		},
	})
	if err != nil {
		return minio.UploadInfo{}, err
	}

	return info, nil
}
