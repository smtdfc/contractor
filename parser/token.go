package parser

import (
	"fmt"
)

type TokenType int

var Keywords map[string]any = map[string]any{
	"model": "",
	"rest":  "",
	"event": "",
}

const (
	TOKEN_NUMBER TokenType = iota
	TOKEN_STRING
	TOKEN_BOOL
	TOKEN_NULL
	TOKEN_IDENTIFIER
	TOKEN_KEYWORD
	TOKEN_OPERATOR
	TOKEN_COMMA
	TOKEN_COLON
	TOKEN_RIGHT_PAREN
	TOKEN_LEFT_PAREN
	TOKEN_RIGHT_BRACE
	TOKEN_LEFT_BRACE
	TOKEN_RIGHT_SQUARE
	TOKEN_LEFT_SQUARE
	TOKEN_ANNOTATION
	TOKEN_NEWLINE
	TOKEN_EOF
)

type Token struct {
	Type  TokenType
	Value string
	Loc   *TokenLocation
}

func (t *Token) ToString() string {
	return fmt.Sprintf("%d : %s", t.Type, t.Value)
}

func (t *Token) HasType(tt TokenType) bool {
	return t.Type == tt
}

func (t *Token) Match(tt TokenType, value string) bool {
	return t.Type == tt && t.Value == value
}

func (t *Token) Copy() *Token {
	return &Token{
		Type:  t.Type,
		Value: t.Value,
		Loc:   t.Loc,
	}
}

type ListToken []*Token

func NewToken(t TokenType, value string, loc *TokenLocation) *Token {
	return &Token{
		Type:  t,
		Value: value,
		Loc:   loc,
	}
}
