package exception

import (
	"errors"
)

type TypeException struct {
	error
	Loc ILocation
	Msg string
}

func (e *TypeException) GetLoc() ILocation {
	return e.Loc
}

func (e *TypeException) GetMsg() string {
	return e.Msg
}

func NewTypeException(msg string, loc ILocation) *TypeException {
	return &TypeException{
		error: errors.New(msg),
		Msg:   msg,
		Loc:   loc,
	}
}
