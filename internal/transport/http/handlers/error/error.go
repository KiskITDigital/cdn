package error

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-faster/jx"
	"github.com/ogen-go/ogen/ogenerrors"
	"github.com/ogen-go/ogen/validate"
	api "gitlab.ubrato.ru/ubrato/cdn/api/gen"
	"gitlab.ubrato.ru/ubrato/cdn/internal/lib/cerr"
)

func (h *Handler) HandleError(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
	var (
		cError                 *cerr.Error
		ogenSecurityError      *ogenerrors.SecurityError
		ogenDecodeParamsError  *ogenerrors.DecodeParamsError
		ogenDecodeRequestError *ogenerrors.DecodeRequestError
	)

	switch {
	case errors.As(err, &cError):
	case errors.As(err, &ogenSecurityError):
		cError = convertOgenSecurityError(ogenSecurityError)
	case errors.As(err, &ogenDecodeParamsError):
		cError = convertOgenDecodeParamsError(ogenDecodeParamsError)
	case errors.As(err, &ogenDecodeRequestError):
		cError = convertOgenDecodeRequestError(ogenDecodeRequestError)
	default:
		cError = cerr.Default(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(cError.Code.HTTPCode())

	apiError, err := newApiError(cError)
	if err != nil {
		h.logger.Error("Create api error", "error", err)
		return
	}

	b, err := json.Marshal(apiError)
	if err != nil {
		h.logger.Error("Marshal api error", "error", err)
		return
	}

	_, err = w.Write(b)
	if err != nil {
		h.logger.Error("Write api error", "error", err)
		return
	}
}

func convertOgenSecurityError(err *ogenerrors.SecurityError) *cerr.Error {
	return cerr.Default(err)
}

func convertOgenDecodeParamsError(err *ogenerrors.DecodeParamsError) *cerr.Error {
	return cerr.Default(err)
}

func convertOgenDecodeRequestError(err *ogenerrors.DecodeRequestError) *cerr.Error {
	var (
		validateError *validate.Error
	)

	switch {
	case errors.As(err, &validateError):
		validationErrors := make(map[string]string)

		for _, field := range validateError.Fields {
			validationErrors[field.Name] = field.Error.Error()
		}

		return cerr.Wrap(err, cerr.CodeValidate, "Validation failed", map[string]interface{}{
			"validation_errors": validationErrors,
		})
	default:
		return cerr.Wrap(err, cerr.CodeValidate, "Unprocessable body", nil)
	}
}

func newApiError(cErr *cerr.Error) (api.WrappedError, error) {
	output := make(map[string]jx.Raw)

	for key, value := range cErr.Attributes {
		data, err := json.Marshal(value)
		if err != nil {
			return api.WrappedError{}, fmt.Errorf("marshal attribute: %w", err)
		}

		output[key] = jx.Raw(data)
	}

	return api.WrappedError{
		Error: api.Error{
			Code:    cErr.Code.String(),
			Message: cErr.Message,
			Details: output,
		},
	}, nil
}
