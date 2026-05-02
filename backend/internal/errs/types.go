package errs

import (
	"net/http"
)

func NewUnauthorizedError(message string, override bool) *HTTPError {
	return &HTTPError{
		Code: MakeUpperCaseWithUnderscores(http.StatusText(http.StatusUnauthorized)),
		Message: message,
		Status: http.StatusUnauthorized,
		Override: override,
	}
}

func NewBadRequestError(message string, override bool, code *string, error []FieldError, action *Action) *HTTPError {
	formattedCode := MakeUpperCaseWithUnderscores(http.StatusText(http.StatusBadRequest))

	if code != nil {
		formattedCode = *code
	}

	return &HTTPError{
		Code: formattedCode,
		Message: message,
		Status: http.StatusNotFound,
		Override: override,
	}
} 

func NewInternalServerError() *HTTPError {
	return &HTTPError{
		Code:     MakeUpperCaseWithUnderscores(http.StatusText(http.StatusInternalServerError)),
		Message:  http.StatusText(http.StatusInternalServerError),
		Status:   http.StatusInternalServerError,
		Override: false,
	}
}

func ValidationError(err error) *HTTPError {
	return NewBadRequestError("Validation failed: "+err.Error(), false, nil, nil, nil)
}