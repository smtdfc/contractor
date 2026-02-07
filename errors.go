package main

import (
	"errors"
)

type BaseError interface {
	GetName() string
	GetMessage() string
	GetLocation() *ErrorLocation
}

type InvalidCharacterError struct {
	error
	Location *ErrorLocation
}

func (e *InvalidCharacterError) GetName() string {
	return "InvalidCharacterError"
}

func (e *InvalidCharacterError) GetMessage() string {
	return e.error.Error()
}

func (e *InvalidCharacterError) GetLocation() *ErrorLocation {
	return e.Location
}

func NewInvalidCharacterError(message string, loc *ErrorLocation) *InvalidCharacterError {
	return &InvalidCharacterError{
		error:    errors.New(message),
		Location: loc,
	}
}

type SyntaxError struct {
	error
	Location *ErrorLocation
}

func (e *SyntaxError) GetName() string {
	return "SyntaxError"
}

func (e *SyntaxError) GetMessage() string {
	return e.error.Error()
}

func (e *SyntaxError) GetLocation() *ErrorLocation {
	return e.Location
}

func NewSyntaxError(message string, loc *ErrorLocation) *SyntaxError {
	return &SyntaxError{
		error:    errors.New(message),
		Location: loc,
	}
}

type TypeError struct {
	error
	Location *ErrorLocation
}

func (e *TypeError) GetName() string {
	return "TypeError"
}

func (e *TypeError) GetMessage() string {
	return e.error.Error()
}

func (e *TypeError) GetLocation() *ErrorLocation {
	return e.Location
}

func NewTypeError(message string, loc *ErrorLocation) *TypeError {
	return &TypeError{
		error:    errors.New(message),
		Location: loc,
	}
}
