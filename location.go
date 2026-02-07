package main

import (
	"fmt"
)

type TokenLocation struct {
	Start Position
	End   Position
	File  string
}

func (l *TokenLocation) Copy() *TokenLocation {
	return &TokenLocation{
		Start: l.Start,
		End:   l.End,
		File:  l.File,
	}
}

func (l *TokenLocation) ToString() string {
	return fmt.Sprintf("%d:%d - %d:%d in %s", l.Start.Line, l.Start.Column, l.End.Line, l.End.Column, l.File)
}

func NewTokenLocation(start Position, end Position, file string) *TokenLocation {
	return &TokenLocation{
		Start: start,
		End:   end,
		File:  file,
	}
}

type NodeLocation struct {
	Start *TokenLocation
	End   *TokenLocation
}

func NewNodeLocation(start *TokenLocation, end *TokenLocation) *NodeLocation {
	return &NodeLocation{
		Start: start,
		End:   end,
	}
}

func (l *NodeLocation) Copy() *NodeLocation {
	return &NodeLocation{
		Start: l.Start.Copy(),
		End:   l.End.Copy(),
	}
}

func (l *NodeLocation) GetErrorLocation() *ErrorLocation {
	return NewErrorLocation(
		l.Start.Start,
		l.End.End,
		l.Start.File,
	)
}

type ErrorLocation struct {
	Start Position
	End   Position
	File  string
}

func (l *ErrorLocation) ToString() string {
	return fmt.Sprintf("%d:%d - %d:%d in %s", l.Start.Line, l.Start.Column, l.End.Line, l.End.Column, l.File)
}

func NewErrorLocation(start Position, end Position, file string) *ErrorLocation {
	return &ErrorLocation{
		Start: start,
		End:   end,
		File:  file,
	}
}
