package helpers

import (
	"fmt"
	"strings"
)

type CodeBuffer struct {
	strings.Builder
	indentSize  int
	indentLevel int
}

func NewCodeBuffer(indentSize int) *CodeBuffer {
	return &CodeBuffer{indentSize: indentSize}
}

func (cb *CodeBuffer) Indent() { cb.indentLevel++ }
func (cb *CodeBuffer) Outdent() {
	if cb.indentLevel > 0 {
		cb.indentLevel--
	}
}

func (cb *CodeBuffer) WriteLine(format string, args ...interface{}) {
	indent := strings.Repeat(" ", cb.indentLevel*cb.indentSize)
	cb.WriteString(indent)
	fmt.Fprintf(&cb.Builder, format, args...)
	cb.WriteString("\n")
}
