package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScanner_Navigation(t *testing.T) {
	input := "A\nB\rC"
	s := NewScanner(input, "source.txt")

	// 1. First character 'A'
	assert.Equal(t, 'A', s.Next())
	assert.Equal(t, 1, s.Line)
	assert.Equal(t, 1, s.Column)

	// 2. Newline '\n' - Still belongs to Line 1
	assert.Equal(t, '\n', s.Next())
	assert.Equal(t, 1, s.Line)
	assert.Equal(t, 2, s.Column)

	// 3. Character 'B' - Now on Line 2
	assert.Equal(t, 'B', s.Next())
	assert.Equal(t, 2, s.Line)
	assert.Equal(t, 1, s.Column)

	// 4. Carriage Return '\r' - Column increments but width is often ignored later
	assert.Equal(t, '\r', s.Next())
	assert.Equal(t, 2, s.Line)
	assert.Equal(t, 2, s.Column)

	// 5. Character 'C' - Column does NOT increment because of previous \r
	assert.Equal(t, 'C', s.Next())
	assert.Equal(t, 2, s.Line)
	assert.Equal(t, 2, s.Column)

	// 6. EOF handling
	assert.Equal(t, rune(0), s.Next())
}

func TestPositionMarker_Capture(t *testing.T) {
	code := "abc"
	scanner := NewScanner(code, "test.go")
	marker := scanner.CreateMarker()

	// Position at 'a'
	scanner.Next()
	marker.MarkStart()

	// Position at 'c'
	scanner.Next() // 'b'
	scanner.Next() // 'c'
	marker.MarkEnd()

	loc, err := marker.GetLocation()
	assert.NoError(t, err)
	assert.Equal(t, 1, loc.Start.Column)
	assert.Equal(t, 3, loc.End.Column)
	assert.Equal(t, "test.go", loc.File)
}
