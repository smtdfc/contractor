package golang

import "errors"

type ValidationError struct {
	error
	Errors map[string]string
}

func (e *ValidationError) Error() string {
	return "validation error"
}

func NewValidationError(d map[string]string) *ValidationError {
	return &ValidationError{
		error:  errors.New("validation error"),
		Errors: d,
	}
}
