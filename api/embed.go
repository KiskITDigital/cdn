package api

import (
	_ "embed"
)

//go:embed bundle.yaml
var OpenapiSpec []byte
