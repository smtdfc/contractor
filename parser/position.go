package parser

import "fmt"

type Position struct {
	Col  int
	Line int
}

func (p *Position) Copy() *Position {
	return &Position{
		Col:  p.Col,
		Line: p.Line,
	}
}

func (p *Position) String() string {
	return fmt.Sprintf("(%d,%d)", p.Line, p.Col)
}

func NewPosition(line int, col int) *Position {
	return &Position{Col: col, Line: line}
}
