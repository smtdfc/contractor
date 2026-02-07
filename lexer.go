package main

type Lexer struct{}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isAlpha(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_'
}

func isAlphaNumeric(r rune) bool {
	return isAlpha(r) || (r >= '0' && r <= '9')
}

func (l *Lexer) getOperatorToken(scanner *Scanner) *Token {
	marker := scanner.CreateMarker()
	marker.MarkStart()

	first := scanner.Current
	scanner.Next()
	literal := string(first)

	if (first == '=' || first == '!' || first == '<' || first == '>' || first == '+' || first == '-' || first == '&' || first == '|') && scanner.Current == '=' {
		literal += string(scanner.Current)
		scanner.Next()
	} else if first == '&' && scanner.Current == '&' {
		literal += string(scanner.Current)
		scanner.Next()
	} else if first == '|' && scanner.Current == '|' {
		literal += string(scanner.Current)
		scanner.Next()
	}

	marker.MarkEnd()
	loc, _ := marker.GetLocation()
	return NewToken(TOKEN_OPERATOR, literal, loc)
}

func (l *Lexer) SkipComment(scanner *Scanner) {
	for scanner.Current != 0 {
		if scanner.Current == '\n' {
			break
		}
		scanner.Next()
	}
	scanner.Next()
}

func (l *Lexer) GetNumberToken(scanner *Scanner) (*Token, BaseError) {
	marker := scanner.CreateMarker()
	marker.MarkStart()

	startIdx := scanner.Index
	isDot := false

	for scanner.Current != 0 {
		if isDigit(scanner.Current) {
			scanner.Next()
		} else if scanner.Current == '.' {
			if isDot {
				return nil, NewInvalidCharacterError(
					"Invalid character: '"+string(scanner.Current)+"'",
					scanner.GetErrorLocation(),
				)
			}
			isDot = true
			scanner.Next()
		} else {
			break
		}

	}

	marker.MarkEnd()

	literal := string(scanner.Code[startIdx:scanner.Index])
	loc, _ := marker.GetLocation()

	return NewToken(TOKEN_NUMBER, literal, loc), nil
}

func (l *Lexer) GetIdentifierToken(scanner *Scanner) (*Token, BaseError) {
	marker := scanner.CreateMarker()
	marker.MarkStart()

	startIdx := scanner.Index
	for scanner.Current != 0 {
		if isAlphaNumeric(scanner.Current) || scanner.Current == '_' {
			scanner.Next()

		} else {
			break
		}
	}

	marker.MarkEnd()
	literal := string(scanner.Code[startIdx:scanner.Index])

	loc, _ := marker.GetLocation()

	tokenType := TOKEN_IDENTIFIER
	if _, ok := Keywords[literal]; ok {
		tokenType = TOKEN_KEYWORD
	}

	return NewToken(tokenType, literal, loc), nil
}

func (l *Lexer) GetStringToken(scanner *Scanner, allowNewline bool, matcher rune) (*Token, BaseError) {
	marker := scanner.CreateMarker()
	marker.MarkStart()

	scanner.Next()
	startIdx := scanner.Index
	for scanner.Current != 0 {
		if scanner.Current == '\n' && !allowNewline {
			return nil, NewInvalidCharacterError(
				"Invalid character: '"+string(scanner.Current)+"'",
				scanner.GetErrorLocation(),
			)
		} else if scanner.Current == matcher {
			break
		} else {
			scanner.Next()
		}
	}

	if scanner.Current == 0 {
		return nil, NewInvalidCharacterError(
			"EOF when scanning string",
			scanner.GetErrorLocation(),
		)
	}

	marker.MarkEnd()
	scanner.Next()

	literal := string(scanner.Code[startIdx:scanner.Index])
	loc, _ := marker.GetLocation()

	return NewToken(TOKEN_STRING, literal, loc), nil

}

func (l *Lexer) getSimpleToken(scanner *Scanner, t TokenType) *Token {
	marker := scanner.CreateMarker()
	marker.MarkStart()
	literal := string(scanner.Current)
	marker.MarkEnd()
	scanner.Next()
	loc, _ := marker.GetLocation()
	return NewToken(t, literal, loc)
}

func (l *Lexer) Parse(code string, file string) (ListToken, BaseError) {
	var tokens ListToken
	scanner := NewScanner(code, file)
	scanner.Next()

	for scanner.Current != 0 {
		switch scanner.Current {
		case '\t', '\r', ' ':
			scanner.Next()
			continue
		case '#':
			l.SkipComment(scanner)
			continue
		case '\n':
			tokens = append(tokens, l.getSimpleToken(scanner, TOKEN_NEWLINE))
			continue
		case '(':
			tokens = append(tokens, l.getSimpleToken(scanner, TOKEN_LEFT_PAREN))
			continue
		case ')':
			tokens = append(tokens, l.getSimpleToken(scanner, TOKEN_RIGHT_PAREN))
			continue
		case '{':
			tokens = append(tokens, l.getSimpleToken(scanner, TOKEN_LEFT_BRACE))
			continue
		case '}':
			tokens = append(tokens, l.getSimpleToken(scanner, TOKEN_RIGHT_BRACE))
			continue
		case '[':
			tokens = append(tokens, l.getSimpleToken(scanner, TOKEN_LEFT_SQUARE))
			continue
		case ']':
			tokens = append(tokens, l.getSimpleToken(scanner, TOKEN_RIGHT_SQUARE))
			continue
		case ':':
			tokens = append(tokens, l.getSimpleToken(scanner, TOKEN_COLON))
			continue
		case ',':
			tokens = append(tokens, l.getSimpleToken(scanner, TOKEN_COMMA))
			continue
		case '@':
			tokens = append(tokens, l.getSimpleToken(scanner, TOKEN_ANNOTATION))
			continue
		case '+', '-', '*', '/', '%', '^', '=', '!', '<', '>', '&', '|':
			tok := l.getOperatorToken(scanner)
			tokens = append(tokens, tok)
			continue

		}

		if isDigit(scanner.Current) {
			tok, err := l.GetNumberToken(scanner)
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, tok)
			continue
		}

		if isAlpha(scanner.Current) || scanner.Current == '_' {
			tok, err := l.GetIdentifierToken(scanner)
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, tok)
			continue
		}

		if scanner.Current == '`' || scanner.Current == '"' {
			char := scanner.Current
			allowNewline := (char == '`')

			tok, err := l.GetStringToken(scanner, allowNewline, char)
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, tok)
			continue
		}

		return nil, NewInvalidCharacterError(
			"Invalid character: '"+string(scanner.Current)+"'",
			scanner.GetErrorLocation(),
		)
	}

	tokens = append(tokens, NewToken(TOKEN_EOF, "", scanner.GetLocation()))
	return tokens, nil
}

func NewLexer() *Lexer {
	return &Lexer{}
}
