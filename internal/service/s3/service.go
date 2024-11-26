package s3

import (
	"github.com/bwmarrin/snowflake"
	"github.com/minio/minio-go/v7"
)

type Service struct {
	minio *minio.Client
	node  *snowflake.Node
}

func New(
	minio *minio.Client,
	node *snowflake.Node,
) *Service {
	return &Service{
		minio: minio,
		node:  node,
	}
}
