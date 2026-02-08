package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestHelper: Creates a scanner and moves to first character
func setupLexerTest(code string) (*Lexer, *Scanner) {
	l := NewLexer()
	s := NewScanner(code, "test_file.txt")
	s.Next() // Initialize first character
	return l, s
}

func TestLexer_GetNumberToken(t *testing.T) {
	// Scenario 1: Valid integer
	l, s := setupLexerTest("12345")
	tok, err := l.GetNumberToken(s)
	assert.NoError(t, err)
	assert.Equal(t, TOKEN_NUMBER, tok.Type)
	assert.Equal(t, "12345", tok.Value)

	// Scenario 2: Valid float
	l, s = setupLexerTest("3.1415")
	tok, err = l.GetNumberToken(s)
	assert.NoError(t, err)
	assert.Equal(t, "3.1415", tok.Value)

	// Scenario 3: Error on multiple dots (Invalid float)
	l, s = setupLexerTest("1.2.3")
	tok, err = l.GetNumberToken(s)
	assert.Error(t, err)
	assert.Nil(t, tok)
}

func TestLexer_GetOperatorToken(t *testing.T) {
	// Scenario 1: Compound assignment operator
	l, s := setupLexerTest("+")
	tok := l.getOperatorToken(s)
	assert.Equal(t, TOKEN_OPERATOR, tok.Type)
	assert.Equal(t, "+", tok.Value)

	// Scenario 2: Logical AND
	l, s = setupLexerTest("&&")
	tok = l.getOperatorToken(s)
	assert.Equal(t, "&&", tok.Value)

	// Scenario 3: Single operator
	l, s = setupLexerTest("*")
	tok = l.getOperatorToken(s)
	assert.Equal(t, "*", tok.Value)
}

func TestLexer_GetIdentifierToken(t *testing.T) {
	// Assuming "if" is in your Keywords map

	// Scenario 1: Normal identifier
	l, s := setupLexerTest("my_variable_123")
	tok, err := l.GetIdentifierToken(s)
	assert.NoError(t, err)
	assert.Equal(t, TOKEN_IDENTIFIER, tok.Type)
	assert.Equal(t, "my_variable_123", tok.Value)

	// Scenario 2: Reserved keyword
	l, s = setupLexerTest("model")
	tok, err = l.GetIdentifierToken(s)
	assert.NoError(t, err)
	assert.Equal(t, TOKEN_KEYWORD, tok.Type)
}

func TestLexer_GetStringToken(t *testing.T) {
	// Scenario 1: Double quoted string (No newline allowed)
	l, s := setupLexerTest("\"hello world\"")
	tok, err := l.GetStringToken(s, false, '"')
	assert.NoError(t, err)

	assert.Equal(t, "hello world", tok.Value)

	// Scenario 2: Backtick string (Newline allowed)
	l, s = setupLexerTest("`line one\nline two` ")
	tok, err = l.GetStringToken(s, true, '`')
	assert.NoError(t, err)
	assert.Equal(t, "line one\nline two", tok.Value)

	// Scenario 3: Error on Unterminated string
	l, s = setupLexerTest("\"unclosed")
	tok, err = l.GetStringToken(s, false, '"')
	assert.Error(t, err, "Should fail when EOF is reached before closing matcher")
}

func TestLexer_Parse(t *testing.T) {
	code := `
		# This is a comment
		var x = 10.5
		if (x > 10) { "Win" }
	`
	l := NewLexer()
	tokens, err := l.Parse(code, "main.src")
	assert.NoError(t, err)

	// Verify sequence of important tokens
	// Skipping newlines/spaces based on your Parse logic
	foundIdentifier := false
	for _, tok := range tokens {
		if tok.Value == "x" {
			foundIdentifier = true
		}
	}
	assert.True(t, foundIdentifier)
	assert.Equal(t, TOKEN_EOF, tokens[len(tokens)-1].Type)
}
