package parser

import "fmt"

type Location struct {
	File  string
	Start *Position
	End   *Position
}

func (l *Location) Copy() *Location {
	return &Location{
		File:  l.File,
		Start: l.Start.Copy(),
		End:   l.End.Copy(),
	}
}

func (l *Location) String() string {
	return fmt.Sprintf("(File: %s ,Start: %s ,End: %s)", l.File, l.Start, l.End)
}

func (l *Location) GetStart() (int, int) {
	return l.Start.Line, l.Start.Col
}

func (l *Location) GetEnd() (int, int) {
	return l.End.Line, l.End.Col
}

func (l *Location) GetFile() string {
	return l.File
}

func NewLocation(file string, start *Position, end *Position) *Location {
	return &Location{
		File:  file,
		Start: start.Copy(),
		End:   end.Copy(),
	}
}
