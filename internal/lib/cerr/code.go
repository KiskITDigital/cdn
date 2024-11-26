package cerr

import "net/http"

type Code string

func (c Code) String() string { return string(c) }

func (c Code) HTTPCode() int {
	httpCode, exist := codes[c]
	if !exist {
		return http.StatusInternalServerError
	}

	return httpCode
}

var (
	CodeInternal           Code = "ERR_INTERNAL"
	CodeUnauthorized       Code = "ERR_UNAUTHORIZED"
	CodeNotPermitted       Code = "ERR_NOT_PERMITTED"
	CodeValidate           Code = "ERR_VALIDATE"
	CodeNotFound           Code = "ERR_NOT_FOUND"
	CodeInvalidCredentials Code = "ERR_INVALID_CREDENTIALS"
)

var codes = map[Code]int{
	CodeInternal:           http.StatusInternalServerError,
	CodeValidate:           http.StatusBadRequest,
	CodeInvalidCredentials: http.StatusBadRequest,
	CodeUnauthorized:       http.StatusUnauthorized,
	CodeNotPermitted:       http.StatusForbidden,
}
