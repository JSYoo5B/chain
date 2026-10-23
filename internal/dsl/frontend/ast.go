package frontend

import (
	"go/ast"
	"go/token"
)

// File keeps the Go package/import header and the DSL workflow declarations.
type File struct {
	Source     *Source
	Header     *ast.File
	HeaderSpan Span
	Workflows  []*Workflow

	headerFile *token.File
}

// HeaderPosition maps a package/import AST position to the original source.
func (f *File) HeaderPosition(pos token.Pos) token.Position {
	if !pos.IsValid() {
		return token.Position{}
	}
	return f.Source.Position(f.headerFile.Offset(pos))
}

// Identifier is a DSL name at its original source location.
type Identifier struct {
	Name string
	Span Span
}

// Workflow contains raw Go regions, their parsed Go AST, and DSL connections.
type Workflow struct {
	Span         Span
	Constructor  Identifier
	Name         Identifier
	Parameters   GoRegion
	StateType    GoRegion
	Prerequisite *GoRegion
	Go           *GoSyntax
	Nodes        []Identifier
	Connections  []Connection
}

// Connection is one directed edge. Direction is success, failure, abort, or a
// decoded branch string. Reverse arrows are normalized to their actual flow.
type Connection struct {
	From      Identifier
	To        Identifier
	Direction string
	Span      Span
}
