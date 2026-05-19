package exception

import (
	"errors"
)

type EmitException struct {
	error
	Loc ILocation
	Msg string
}

func (e *EmitException) GetLoc() ILocation {
	return e.Loc
}

func (e *EmitException) GetMsg() string {
	return e.Msg
}

func NewEmitException(msg string, loc ILocation) *EmitException {
	return &EmitException{
		error: errors.New(msg),
		Msg:   msg,
		Loc:   loc,
	}
}
