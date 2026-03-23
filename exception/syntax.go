package exception

import (
	"errors"
)

type SyntaxException struct {
	error
	Loc ILocation
	Msg string
}

func (e *SyntaxException) GetLoc() ILocation {
	return e.Loc
}

func (e *SyntaxException) GetMsg() string {
	return e.Msg
}

func NewSyntaxException(msg string, loc ILocation) *SyntaxException {
	return &SyntaxException{
		error: errors.New(msg),
		Msg:   msg,
		Loc:   loc,
	}
}
