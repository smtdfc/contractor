package parser

type Position struct {
	Column int
	Line   int
}

func NewPosition(line int, column int) Position {
	return Position{
		Line:   line,
		Column: column,
	}
}
