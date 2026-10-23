package frontend

import (
	"fmt"
	"go/scanner"
	"go/token"
)

// A lexeme records Go's tokenization without classifying DSL keywords.
// literal is the Go scanner's value, not necessarily the original source text:
// for example, Go removes carriage returns from raw string literals.
type lexeme struct {
	kind    token.Token
	offset  int
	literal string
}

type cursor struct {
	source *Source
	tokens []lexeme
	index  int
}

func scan(source *Source) (*cursor, error) {
	var errors scanner.ErrorList
	var lexer scanner.Scanner
	// Keep scanner line directives separate from the source's physical positions.
	file := token.NewFileSet().AddFile(source.file.Name(), -1, len(source.text))
	lexer.Init(file, source.text, func(pos token.Position, message string) {
		errors.Add(source.Position(pos.Offset), message)
	}, scanner.ScanComments)

	c := &cursor{source: source}
	for {
		pos, kind, literal := lexer.Scan()
		if kind == token.COMMENT {
			continue
		}
		c.tokens = append(c.tokens, lexeme{
			kind: kind, offset: file.Offset(pos), literal: literal,
		})
		if kind == token.EOF {
			break
		}
	}
	if err := errors.Err(); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *cursor) current() lexeme {
	return c.tokens[c.index]
}

func (c *cursor) advance() {
	if c.current().kind != token.EOF {
		c.index++
	}
}

// GoRegion preserves both the delimiters and the source enclosed by them.
// Extracting a region does not validate the Go syntax inside it.
type GoRegion struct {
	Span    Span
	Content Span
}

type delimiter struct {
	open  token.Token
	close token.Token
	start int
}

func closingToken(open token.Token) token.Token {
	switch open {
	case token.LPAREN:
		return token.RPAREN
	case token.LBRACK:
		return token.RBRACK
	case token.LBRACE:
		return token.RBRACE
	default:
		return token.ILLEGAL
	}
}

// takeGoRegion consumes a balanced (), [], or {} region. Strings and comments
// are already tokenized, so their delimiters cannot affect nesting.
func (c *cursor) takeGoRegion() (GoRegion, error) {
	first := c.current()
	if closingToken(first.kind) == token.ILLEGAL {
		return GoRegion{}, &scanner.Error{
			Pos: c.source.Position(first.offset),
			Msg: "expected '(', '[' or '{' to start Go region",
		}
	}

	var stack []delimiter
	for {
		current := c.current()
		if close := closingToken(current.kind); close != token.ILLEGAL {
			stack = append(stack, delimiter{current.kind, close, current.offset})
		} else if current.kind == token.RPAREN || current.kind == token.RBRACK || current.kind == token.RBRACE {
			top := stack[len(stack)-1]
			if current.kind != top.close {
				return GoRegion{}, &scanner.Error{
					Pos: c.source.Position(current.offset),
					Msg: fmt.Sprintf("expected %q to close %q at %s, found %q",
						top.close, top.open, c.source.Position(top.start), current.kind),
				}
			}
			stack = stack[:len(stack)-1]
			if len(stack) == 0 {
				region := GoRegion{
					Span:    Span{Start: first.offset, End: current.offset + 1},
					Content: Span{Start: first.offset + 1, End: current.offset},
				}
				c.advance()
				return region, nil
			}
		} else if current.kind == token.EOF {
			top := stack[len(stack)-1]
			return GoRegion{}, &scanner.Error{
				Pos: c.source.Position(current.offset),
				Msg: fmt.Sprintf("expected %q to close %q at %s, found end of file",
					top.close, top.open, c.source.Position(top.start)),
			}
		}
		c.advance()
	}
}
