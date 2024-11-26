package http

import (
	"context"
	"net/http"

	api "gitlab.ubrato.ru/ubrato/cdn/api/gen"
)

var _ api.Handler = new(Router)

type Router struct {
	Error
	File
	Auth
}

type Error interface {
	HandleError(ctx context.Context, w http.ResponseWriter, r *http.Request, err error)
}

type File interface {
	FileIDGet(ctx context.Context, params api.FileIDGetParams) (api.FileIDGetRes, error)
	FileIDHead(ctx context.Context, params api.FileIDHeadParams) (api.FileIDHeadRes, error)
	UploadPost(ctx context.Context, req *api.UploadPostReq, params api.UploadPostParams) (api.UploadPostRes, error)
}

type Auth interface {
	HandleBearerAuth(ctx context.Context, operationName string, t api.BearerAuth) (context.Context, error)
}

type RouterParams struct {
	Error Error
	File  File
	Auth  Auth
}

func NewRouter(params RouterParams) *Router {
	return &Router{
		Error: params.Error,
		File:  params.File,
		Auth:  params.Auth,
	}
}
