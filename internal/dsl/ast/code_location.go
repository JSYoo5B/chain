package ast

import "fmt"

type CodeLocation struct {
	Line   int
	Column int
	Text   string
}

func (l CodeLocation) Location() string {
	return fmt.Sprintf("%d:%d", l.Line, l.Column)
}
