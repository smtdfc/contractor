package parser

const nullRune = '\x00'

type Scanner struct {
	Current rune
	Code    string
	Index   int
	Col     int
	Line    int
	NextIdx int
}

func (s *Scanner) Next() rune {
	if s.NextIdx >= len(s.Code) {
		s.Current = nullRune
		return nullRune
	}

	s.Current = rune(s.Code[s.NextIdx])
	s.NextIdx++

	if s.Current == '\n' {
		s.Line++
		s.Col = 0
	} else {
		s.Col++
	}

	return s.Current
}

func (s *Scanner) Peek() rune {
	if s.NextIdx >= len(s.Code) {
		return nullRune
	}
	return rune(s.Code[s.NextIdx])
}

func (s *Scanner) GetPosition() *Position {
	return &Position{
		Col:  s.Col,
		Line: s.Line,
	}
}

func (s *Scanner) skipWhitespace() {
	for {
		switch s.Current {
		case ' ', '\t', '\r', '\n':
			s.Next()
		default:
			return
		}
	}
}

func NewScanner(code string) *Scanner {
	return &Scanner{Code: code, Index: -1, Col: 0, Line: 1}
}
