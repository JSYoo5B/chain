// Package frontend parses Chain DSL source while preserving embedded Go code.
package frontend

import (
	"bytes"
	"go/token"
)

// Span is a half-open range of byte offsets in the original source.
type Span struct {
	Start int
	End   int
}

// Source owns an immutable copy of a DSL file and its physical positions.
type Source struct {
	text []byte
	file *token.File
}

// NewSource copies text so later caller edits cannot invalidate source spans.
func NewSource(filename string, text []byte) *Source {
	text = bytes.Clone(text)
	file := token.NewFileSet().AddFile(filename, -1, len(text))
	file.SetLinesForContent(text)
	return &Source{text: text, file: file}
}

// Text returns the original bytes in a span, including comments and whitespace.
func (s *Source) Text(span Span) string {
	return string(s.text[span.Start:span.End])
}

// Position resolves an offset without applying Go //line directives.
func (s *Source) Position(offset int) token.Position {
	return s.file.PositionFor(s.file.Pos(offset), false)
}
