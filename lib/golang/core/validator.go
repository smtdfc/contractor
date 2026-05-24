package core

import (
	"fmt"
	"strings"
)

type ValidateResult[T any] struct {
	Data      *T
	IsSuccess bool
	Errors    map[string]string
}

type IValidator[T any] interface {
	Validate(model *T) error
}

type ValidationErrorDetails map[string]string

// ValidationError is returned by validators when one or more fields fail validation.
// It implements the error interface and exposes a `Details` map keyed by field name.
type ValidationError struct {
	Details ValidationErrorDetails `json:"details"`
}

func (ve ValidationError) Error() string {
	if len(ve.Details) == 0 {
		return "validation failed"
	}

	parts := make([]string, 0, len(ve.Details))
	for k, v := range ve.Details {
		parts = append(parts, fmt.Sprintf("%s: %s", k, v))
	}

	return "validation failed: " + strings.Join(parts, ", ")
}

// NewValidationError returns an empty ValidationError.
func NewValidationError() *ValidationError {
	return &ValidationError{Details: ValidationErrorDetails{}}
}

// NewValidationErrorFromMap constructs a ValidationError from an existing map.
func NewValidationErrorFromMap(m map[string]string) *ValidationError {
	d := ValidationErrorDetails{}
	for k, v := range m {
		d[k] = v
	}
	return &ValidationError{Details: d}
}

// Add adds or replaces a field error on the ValidationError.
func (ve *ValidationError) Add(field, msg string) {
	if ve.Details == nil {
		ve.Details = ValidationErrorDetails{}
	}
	ve.Details[field] = msg
}

// Merge merges another ValidationError's details into this one.
func (ve *ValidationError) Merge(other ValidationError) {
	if other.Details == nil {
		return
	}
	if ve.Details == nil {
		ve.Details = ValidationErrorDetails{}
	}
	for k, v := range other.Details {
		ve.Details[k] = v
	}
}

// HasErrors returns true when there are one or more validation failures.
func (ve ValidationError) HasErrors() bool {
	return len(ve.Details) > 0
}

// IsValidationError reports whether err is a ValidationError.
func IsValidationError(err error) bool {
	if err == nil {
		return false
	}
	switch err.(type) {
	case ValidationError, *ValidationError:
		return true
	default:
		return false
	}
}
