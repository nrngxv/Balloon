package errs

import (

	"strings"
)

type FieldError struct {
	Field string
	Error string
}

type ActionType string

const (
	ActionTypeRedirect ActionType = "redirect"
)

type Action struct {
	Type 	ActionType
	Message string
	Value 	string
}

type HTTPError struct {
	Code		string
	Message		string
	Status		int
	Override	bool

	Error 		[]FieldError
	Action		*Action
}

func (e *HTTPError) Error() string {
	return e.Message
}

func (e *HTTPError) is(target error) bool {
	_, ok := target.(*HTTPError)

	return ok
}

func (e *HTTPError) WithMessage(message string) *HTTPError {
	return &HTTPError{
		Code: e.Code,
		Message: message,
		Status: e.Status,
		Override: e.Override,
		Error: e.Error,
		Action: e.Action,
	}
}

func MakeUpperCaseWithUnderscores(str string) string {
	return strings.ToUpper(strings.Replace(str, " ", "_"))
}