package main

type PositionMarker struct {
	Scanner *Scanner
	Start   Position
	End     Position
}

func (m *PositionMarker) MarkStart() {
	m.Start = NewPosition(m.Scanner.Line, m.Scanner.Column)
}

func (m *PositionMarker) MarkEnd() {
	m.End = NewPosition(m.Scanner.Line, m.Scanner.Column)
}

func (m *PositionMarker) GetErrorLocation() *ErrorLocation {
	return NewErrorLocation(
		m.Start,
		m.End,
		m.Scanner.File,
	)
}

func (m *PositionMarker) GetLocation() (*TokenLocation, error) {
	return NewTokenLocation(
		m.Start,
		m.End,
		m.Scanner.File,
	), nil
}

func NewPositionMarker(scanner *Scanner) *PositionMarker {
	return &PositionMarker{
		Scanner: scanner,
	}
}

type Scanner struct {
	Current rune
	Index   int
	Code    []rune
	File    string
	Column  int
	Line    int
}

func (s *Scanner) CreateMarker() *PositionMarker {
	return NewPositionMarker(s)
}

func (s *Scanner) Next() rune {
	if s.Current == '\n' {
		s.Line += 1
		s.Column = 1
	} else if s.Current != 0 && s.Current != '\r' {
		s.Column += 1
	}

	s.Index += 1
	if s.Index >= len(s.Code) {
		s.Current = 0
		return 0
	}

	s.Current = s.Code[s.Index]

	if s.Index == 0 {
		s.Line = 1
		s.Column = 1
	}

	return s.Current
}

func (s *Scanner) GetErrorLocation() *ErrorLocation {
	if s.Index >= len(s.Code) {
		return nil
	}

	return &ErrorLocation{
		Start: NewPosition(s.Line, s.Column),
		End:   NewPosition(s.Line, s.Column),
		File:  s.File,
	}
}

func (s *Scanner) GetLocation() *TokenLocation {

	return &TokenLocation{
		Start: NewPosition(s.Line, s.Column),
		End:   NewPosition(s.Line, s.Column),
		File:  s.File,
	}
}

func NewScanner(code string, file string) *Scanner {
	return &Scanner{
		Code:    []rune(code),
		Current: 0,
		Index:   -1,
		Column:  1,
		Line:    1,
		File:    file,
	}
}

type TokenScanner struct {
	Tokens  ListToken
	Index   int
	Current *Token
}

func (s *TokenScanner) GetLocation() *NodeLocation {
	return NewNodeLocation(
		s.Current.Loc.Copy(),
		s.Current.Loc.Copy(),
	)
}

func (s *TokenScanner) GetErrorLocation() *ErrorLocation {
	return NewErrorLocation(
		s.Current.Loc.Start,
		s.Current.Loc.End,
		s.Current.Loc.File,
	)
}

func (s *TokenScanner) Next() *Token {
	s.Index++
	if s.Index >= len(s.Tokens) {
		s.Current = nil
	} else {
		s.Current = s.Tokens[s.Index]
	}
	return s.Current
}

func NewTokenScanner(list ListToken) *TokenScanner {
	return &TokenScanner{
		Tokens:  list,
		Index:   -1,
		Current: nil,
	}
}
