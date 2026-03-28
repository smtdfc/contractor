package parser

import (
	"fmt"
	"strings"

	"github.com/smtdfc/contractor/exception"
)

type Lexer struct {
	File string
}

func (l *Lexer) ReadNumber(scanner *Scanner) (*Token, exception.IException) {
	startPos := scanner.GetPosition()
	var num strings.Builder

	num.WriteRune(scanner.Current)

	isDotted := false
	lastValidPos := startPos

	scanner.Next()
	for scanner.Current != nullRune {
		if IsDigit(scanner.Current) {
			num.WriteRune(scanner.Current)
			lastValidPos = scanner.GetPosition()
		} else if scanner.Current == '.' && !isDotted {
			peek := scanner.Peek()
			if !IsDigit(peek) {
				break
			}
			isDotted = true
			num.WriteRune(scanner.Current)
			lastValidPos = scanner.GetPosition()
		} else {
			break
		}
		scanner.Next()
	}

	resultStr := num.String()
	if strings.HasSuffix(resultStr, ".") {
		return nil, exception.NewSyntaxException(
			"Numbers cannot end with a dot",
			NewLocation(l.File, lastValidPos, lastValidPos),
		)
	}

	return NewToken(
		TT_NUMBER,
		resultStr,
		NewLocation(l.File, startPos, lastValidPos),
	), nil
}

func (l *Lexer) ReadIdent(scanner *Scanner) (*Token, exception.IException) {
	startPos := scanner.GetPosition()
	var ident strings.Builder

	ident.WriteRune(scanner.Current)

	lastValidPos := startPos

	scanner.Next()
	for scanner.Current != nullRune {
		if IsAlphaNumeric(scanner.Current) {
			ident.WriteRune(scanner.Current)
			lastValidPos = scanner.GetPosition()
		} else {
			break
		}
		scanner.Next()
	}

	resultStr := ident.String()

	return NewToken(
		TT_IDENT,
		resultStr,
		NewLocation(l.File, startPos, lastValidPos),
	), nil
}

func (l *Lexer) ReadString(scanner *Scanner) (*Token, exception.IException) {
	startPos := scanner.GetPosition()
	var value strings.Builder

	scanner.Next()
	lastPos := startPos

	for scanner.Current != nullRune {
		if scanner.Current == '"' {
			endPos := scanner.GetPosition()
			scanner.Next()
			return NewToken(
				TT_STRING,
				value.String(),
				NewLocation(l.File, startPos, endPos),
			), nil
		}

		if scanner.Current == '\\' {
			scanner.Next()
			switch scanner.Current {
			case 'n':
				value.WriteRune('\n')
			case 't':
				value.WriteRune('\t')
			case 'r':
				value.WriteRune('\r')
			case '\\':
				value.WriteRune('\\')
			case '"':
				value.WriteRune('"')
			default:
				return nil, exception.NewSyntaxException(
					fmt.Sprintf("Invalid escape sequence '\\%s'", string(scanner.Current)),
					NewLocation(l.File, scanner.GetPosition(), scanner.GetPosition()),
				)
			}
			lastPos = scanner.GetPosition()
			scanner.Next()
			continue
		}

		if scanner.Current == '\n' {
			return nil, exception.NewSyntaxException(
				"Unterminated string literal",
				NewLocation(l.File, startPos, lastPos),
			)
		}

		value.WriteRune(scanner.Current)
		lastPos = scanner.GetPosition()
		scanner.Next()
	}

	return nil, exception.NewSyntaxException(
		"Unterminated string literal",
		NewLocation(l.File, startPos, lastPos),
	)
}

func (l *Lexer) Start(code string) (TokenList, exception.IException) {
	tokens := make(TokenList, 0)
	scanner := NewScanner(code)
	scanner.Next()
	scanner.skipWhitespace()

	for scanner.Current != nullRune {
		switch {
		case scanner.Current == ' ' || scanner.Current == '\t' || scanner.Current == '\r':
			scanner.Next()

		case scanner.Current == '"':
			token, err := l.ReadString(scanner)
			if err != nil {
				return nil, err
			}

			tokens = append(tokens, token)
			continue

		case IsDigit(scanner.Current):
			token, err := l.ReadNumber(scanner)
			if err != nil {
				return nil, err
			}

			tokens = append(tokens, token)
			continue

		case IsAlpha(scanner.Current):
			token, err := l.ReadIdent(scanner)
			if err != nil {
				return nil, err
			}

			tokens = append(tokens, token)
			continue

		case scanner.Current == '{':
			currentPos := scanner.GetPosition()
			tokens = append(tokens, NewToken(
				TT_LBRACE,
				"{",
				NewLocation(l.File, currentPos, currentPos),
			))
			scanner.Next()
			continue

		case scanner.Current == '}':
			currentPos := scanner.GetPosition()
			tokens = append(tokens, NewToken(
				TT_RBRACE,
				"}",
				NewLocation(l.File, currentPos, currentPos),
			))
			scanner.Next()
			continue

		case scanner.Current == '(':
			currentPos := scanner.GetPosition()
			tokens = append(tokens, NewToken(
				TT_LPAREN,
				"(",
				NewLocation(l.File, currentPos, currentPos),
			))
			scanner.Next()
			continue

		case scanner.Current == ')':
			currentPos := scanner.GetPosition()
			tokens = append(tokens, NewToken(
				TT_RPAREN,
				")",
				NewLocation(l.File, currentPos, currentPos),
			))
			scanner.Next()
			continue

		case scanner.Current == '[':
			currentPos := scanner.GetPosition()
			tokens = append(tokens, NewToken(
				TT_LSQUARE,
				"[",
				NewLocation(l.File, currentPos, currentPos),
			))
			scanner.Next()
			continue

		case scanner.Current == ']':
			currentPos := scanner.GetPosition()
			tokens = append(tokens, NewToken(
				TT_RSQUARE,
				"]",
				NewLocation(l.File, currentPos, currentPos),
			))
			scanner.Next()
			continue

		case scanner.Current == ':':
			currentPos := scanner.GetPosition()
			tokens = append(tokens, NewToken(
				TT_COLON,
				":",
				NewLocation(l.File, currentPos, currentPos),
			))
			scanner.Next()
			continue

		case scanner.Current == '.':
			currentPos := scanner.GetPosition()
			tokens = append(tokens, NewToken(
				TT_DOT,
				".",
				NewLocation(l.File, currentPos, currentPos),
			))
			scanner.Next()
			continue

		case scanner.Current == ',':
			currentPos := scanner.GetPosition()
			tokens = append(tokens, NewToken(
				TT_COMMA,
				",",
				NewLocation(l.File, currentPos, currentPos),
			))
			scanner.Next()
			continue

		case scanner.Current == '@':
			currentPos := scanner.GetPosition()
			tokens = append(tokens, NewToken(
				TT_DECORATOR,
				"@",
				NewLocation(l.File, currentPos, currentPos),
			))
			scanner.Next()
			continue

		case scanner.Current == '?':
			currentPos := scanner.GetPosition()
			tokens = append(tokens, NewToken(
				TT_QUES,
				"?",
				NewLocation(l.File, currentPos, currentPos),
			))
			scanner.Next()
			continue

		case scanner.Current == '<':
			currentPos := scanner.GetPosition()
			tokens = append(tokens, NewToken(
				TT_OP,
				"<",
				NewLocation(l.File, currentPos, currentPos),
			))
			scanner.Next()
			continue

		case scanner.Current == '>':
			currentPos := scanner.GetPosition()
			tokens = append(tokens, NewToken(
				TT_OP,
				">",
				NewLocation(l.File, currentPos, currentPos),
			))
			scanner.Next()
			continue

		case scanner.Current == '\n':
			currentPos := scanner.GetPosition()
			tokens = append(tokens, NewToken(
				TT_NEWLINE,
				"Newline",
				NewLocation(l.File, currentPos, currentPos),
			))
			scanner.Next()
			continue

		default:
			currentPos := scanner.GetPosition()
			return nil, exception.NewSyntaxException(
				fmt.Sprintf("Invalid character '%s' ", string(scanner.Current)),
				NewLocation(
					l.File,
					currentPos,
					currentPos,
				),
			)

		}
	}

	tokens = append(tokens, NewToken(TT_EOF, "", NewLocation(l.File, scanner.GetPosition(), scanner.GetPosition())))

	return tokens, nil

}

func NewLexer(file string) *Lexer {
	return &Lexer{File: file}
}
