package parser

import "fmt"

var count = 0

func ins() TokenType {
	count++
	return TokenType(fmt.Sprintf("%d", count))
}

type TokenType string

var (
	TT_STRING    TokenType = ins()
	TT_NUMBER    TokenType = ins()
	TT_IDENT     TokenType = ins()
	TT_EOF       TokenType = ins()
	TT_LBRACE    TokenType = ins()
	TT_RBRACE    TokenType = ins()
	TT_LPAREN    TokenType = ins()
	TT_RPAREN    TokenType = ins()
	TT_LSQUARE   TokenType = ins()
	TT_RSQUARE   TokenType = ins()
	TT_COLON     TokenType = ins()
	TT_DOT       TokenType = ins()
	TT_COMMA     TokenType = ins()
	TT_DECORATOR TokenType = ins()
	TT_QUES      TokenType = ins()
	TT_NEWLINE   TokenType = ins()
	TT_OP        TokenType = ins()
)

type Token struct {
	Type  TokenType
	Value string
	Loc   *Location
}

func (t *Token) String() string {
	return fmt.Sprintf("Type: %s | Value: %s | Loc: %s", t.Type, t.Value, t.Loc)
}

func (t *Token) Match(t_ TokenType, value string) bool {
	return t.Type == t_ && t.Value == value
}

func (t *Token) MatchType(t_ TokenType) bool {
	return t.Type == t_
}

func NewToken(t TokenType, value string, loc *Location) *Token {
	return &Token{Type: t, Value: value, Loc: loc}
}

type TokenList []*Token
