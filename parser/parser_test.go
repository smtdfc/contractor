package parser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// Helper function to create a basic token for testing
func createTestToken(t TokenType, val string) *Token {
	return &Token{
		Type:  t,
		Value: val,
		Loc:   &TokenLocation{File: "test.src", Start: Position{1, 1}, End: Position{1, 1}},
	}
}

func TestParser_ParseLiteral(t *testing.T) {
	p := NewParser()

	// Test Case: Valid Number Literal
	tokens := ListToken{
		createTestToken(TOKEN_NUMBER, "100"),
		createTestToken(TOKEN_EOF, ""),
	}
	ts := NewTokenScanner(tokens)
	ts.Next()

	node, err := p.ParseLiteral(ts)
	assert.NoError(t, err)
	literal := node.(*LiteralNode)
	assert.Equal(t, "100", literal.Value)

	// Test Case: Invalid Literal (Identifier)
	tokens = ListToken{createTestToken(TOKEN_IDENTIFIER, "wrong")}
	ts = NewTokenScanner(tokens)
	ts.Next()
	_, err = p.ParseLiteral(ts)
	assert.Error(t, err, "Should fail when token is not a literal type")
}

func TestParser_ParseType(t *testing.T) {
	p := NewParser()

	// Scenario: Complex type with generics: List<String>
	tokens := ListToken{
		createTestToken(TOKEN_IDENTIFIER, "List"),
		createTestToken(TOKEN_OPERATOR, "<"),
		createTestToken(TOKEN_IDENTIFIER, "String"),
		createTestToken(TOKEN_OPERATOR, ">"),
		createTestToken(TOKEN_EOF, ""),
	}
	ts := NewTokenScanner(tokens)
	ts.Next()

	typeNode, err := p.ParseType(ts)
	assert.NoError(t, err)
	assert.Equal(t, "List", typeNode.Name)

	// Check generic inner type
	innerType := typeNode.Generic.(*TypeDeclarationNode)
	assert.Equal(t, "String", innerType.Name)
}

func TestParser_ParseAnnotation(t *testing.T) {
	p := NewParser()

	// Scenario: Annotation with arguments @Table("users", true)
	tokens := ListToken{
		createTestToken(TOKEN_ANNOTATION, "@"),
		createTestToken(TOKEN_IDENTIFIER, "Table"),
		createTestToken(TOKEN_LEFT_PAREN, "("),
		createTestToken(TOKEN_STRING, "users"),
		createTestToken(TOKEN_COMMA, ","),
		createTestToken(TOKEN_BOOL, "true"),
		createTestToken(TOKEN_RIGHT_PAREN, ")"),
		createTestToken(TOKEN_NEWLINE, "\n"),
		createTestToken(TOKEN_EOF, ""),
	}
	ts := NewTokenScanner(tokens)
	ts.Next()

	anno, err := p.ParseAnnotation(ts)
	assert.NoError(t, err)
	assert.Equal(t, "Table", anno.Name)
	assert.Len(t, anno.Args, 2)
}

func TestParser_ParseModelStatement(t *testing.T) {
	p := NewParser()

	// Scenario: model User { String name }
	tokens := ListToken{
		createTestToken(TOKEN_KEYWORD, "model"),
		createTestToken(TOKEN_IDENTIFIER, "User"),
		createTestToken(TOKEN_LEFT_BRACE, "{"),
		createTestToken(TOKEN_IDENTIFIER, "String"),
		createTestToken(TOKEN_IDENTIFIER, "name"),
		createTestToken(TOKEN_RIGHT_BRACE, "}"),
		createTestToken(TOKEN_EOF, ""),
	}
	ts := NewTokenScanner(tokens)
	ts.Next()

	model, err := p.ParseModelStatement(ts, nil)
	assert.NoError(t, err)
	assert.Equal(t, "User", model.Name)
	assert.Len(t, model.Fields, 1)
	assert.Equal(t, "name", model.Fields[0].Name)
}
