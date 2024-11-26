package cerr

import "fmt"

type Error struct {
	Code       Code
	Message    string
	Attributes map[string]interface{}
	err        error
}

func Wrap(err error, code Code, msg string, attributes map[string]interface{}) *Error {
	if attributes == nil {
		attributes = make(map[string]interface{})
	}

	attributes["error_trace"] = err.Error()

	return &Error{
		Code:       code,
		Message:    msg,
		Attributes: attributes,
		err:        err,
	}
}

func Default(err error) *Error {
	return &Error{
		Code:    CodeInternal,
		Message: "Unexpected error",
		Attributes: map[string]interface{}{
			"error_trace": err.Error(),
		},
		err: err,
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s (%s): %s", e.Message, e.Code, e.err)
}

func (e *Error) Unwrap() error {
	return e.err
}
