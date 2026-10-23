package frontend

import (
	"fmt"
	goparser "go/parser"
	"go/scanner"
	"go/token"
	"strconv"
)

// Parse reads a DSL file, validates embedded Go syntax and prerequisite return
// scope, and requires a nonempty node list and at least one connection per
// workflow. Name resolution, graph validation, and Go type checking are later
// compiler stages. Any error returns a nil file.
func Parse(filename string, text []byte) (*File, error) {
	source := NewSource(filename, text)
	c, err := scan(source)
	if err != nil {
		return nil, err
	}
	p := &dslParser{cursor: c}
	file, err := p.header()
	if err != nil {
		return nil, err
	}
	for p.current().kind != token.EOF {
		workflow, err := p.workflow()
		if err != nil {
			return nil, err
		}
		file.Workflows = append(file.Workflows, workflow)
		if err := p.endStatement(); err != nil {
			return nil, err
		}
	}
	return file, nil
}

type dslParser struct{ *cursor }

func (p *dslParser) fail(message string) error {
	return &scanner.Error{Pos: p.source.Position(p.current().offset), Msg: message}
}

func (p *dslParser) expect(kind token.Token) (lexeme, error) {
	current := p.current()
	if current.kind != kind {
		return lexeme{}, p.fail(fmt.Sprintf("expected %q, found %q", kind, current.kind))
	}
	p.advance()
	return current, nil
}

func (p *dslParser) identifier() (Identifier, error) {
	current, err := p.expect(token.IDENT)
	if err != nil {
		return Identifier{}, err
	}
	return Identifier{Name: current.literal, Span: Span{Start: current.offset, End: current.offset + len(current.literal)}}, nil
}

func (p *dslParser) skipSemicolons() {
	for p.current().kind == token.SEMICOLON {
		p.advance()
	}
}

// Soft keywords are Go identifiers, so discard the automatic semicolon Go's
// scanner may insert after one. Explicit semicolons remain syntax separators.
func (p *dslParser) skipImplicitSemicolons() {
	for p.current().kind == token.SEMICOLON && p.current().literal == "\n" {
		p.advance()
	}
}

func (p *dslParser) word(word string) error {
	if p.current().kind != token.IDENT || p.current().literal != word {
		return p.fail(fmt.Sprintf("expected %q", word))
	}
	p.advance()
	p.skipImplicitSemicolons()
	return nil
}

func (p *dslParser) endStatement() error {
	if p.current().kind == token.SEMICOLON {
		p.skipSemicolons()
		return nil
	}
	if p.current().kind == token.RBRACE || p.current().kind == token.EOF {
		return nil
	}
	return p.fail("expected newline or ';' after statement")
}

func (p *dslParser) region(open token.Token) (GoRegion, error) {
	if p.current().kind != open {
		return GoRegion{}, p.fail(fmt.Sprintf("expected %q to start Go region", open))
	}
	return p.takeGoRegion()
}

func (p *dslParser) header() (*File, error) {
	if _, err := p.expect(token.PACKAGE); err != nil {
		return nil, err
	}
	if _, err := p.identifier(); err != nil {
		return nil, err
	}
	if err := p.endStatement(); err != nil {
		return nil, err
	}
	for p.current().kind == token.IMPORT {
		p.advance()
		if p.current().kind == token.LPAREN {
			if _, err := p.takeGoRegion(); err != nil {
				return nil, err
			}
		} else {
			if p.current().kind == token.IDENT || p.current().kind == token.PERIOD {
				p.advance()
			}
			if _, err := p.expect(token.STRING); err != nil {
				return nil, err
			}
		}
		if err := p.endStatement(); err != nil {
			return nil, err
		}
	}
	end := p.current().offset
	fset := token.NewFileSet()
	header, err := goparser.ParseFile(fset, p.source.file.Name(), p.source.text[:end],
		goparser.AllErrors|goparser.ParseComments|goparser.SkipObjectResolution)
	if err != nil {
		return nil, remapGoErrors(err, p.source.Position)
	}
	return &File{
		Source: p.source, Header: header, HeaderSpan: Span{Start: 0, End: end},
		headerFile: fset.File(header.Package),
	}, nil
}

func (p *dslParser) workflow() (*Workflow, error) {
	start := p.current().offset
	if err := p.word("workflow"); err != nil {
		return nil, err
	}
	w := &Workflow{}
	var err error
	if w.Constructor, err = p.identifier(); err != nil {
		return nil, err
	}
	if w.Parameters, err = p.region(token.LPAREN); err != nil {
		return nil, err
	}
	if w.Name, err = p.identifier(); err != nil {
		return nil, err
	}
	if w.StateType, err = p.region(token.LBRACK); err != nil {
		return nil, err
	}
	if _, err := p.expect(token.LBRACE); err != nil {
		return nil, err
	}
	p.skipSemicolons()
	if p.current().kind == token.IDENT && p.current().literal == "prerequisite" {
		if err := p.word("prerequisite"); err != nil {
			return nil, err
		}
		region, err := p.region(token.LBRACE)
		if err != nil {
			return nil, err
		}
		w.Prerequisite = &region
		p.skipSemicolons()
	}
	if err := p.word("nodes"); err != nil {
		return nil, err
	}
	if _, err := p.expect(token.COLON); err != nil {
		return nil, err
	}
	if p.labelAhead() || p.current().kind == token.RBRACE || p.current().kind == token.EOF {
		return nil, p.fail("nodes section requires at least one node")
	}
	for {
		node, err := p.identifier()
		if err != nil {
			return nil, err
		}
		w.Nodes = append(w.Nodes, node)
		if p.current().kind != token.COMMA {
			break
		}
		p.advance()
	}
	if err := p.endStatement(); err != nil {
		return nil, err
	}
	for p.current().kind != token.RBRACE && p.current().kind != token.EOF {
		section, err := p.identifier()
		if err != nil {
			return nil, err
		}
		switch section.Name {
		case "success", "failure", "abort", "branches":
		default:
			return nil, &scanner.Error{Pos: p.source.Position(section.Span.Start), Msg: "unknown connection section " + strconv.Quote(section.Name)}
		}
		p.skipImplicitSemicolons()
		if _, err := p.expect(token.COLON); err != nil {
			return nil, err
		}
		for p.current().kind != token.RBRACE && p.current().kind != token.EOF && !p.labelAhead() {
			if section.Name == "branches" {
				connection, err := p.branch()
				if err != nil {
					return nil, err
				}
				w.Connections = append(w.Connections, connection)
			} else {
				connections, err := p.directions(section.Name)
				if err != nil {
					return nil, err
				}
				w.Connections = append(w.Connections, connections...)
			}
			if err := p.endStatement(); err != nil {
				return nil, err
			}
		}
	}
	if len(w.Connections) == 0 {
		return nil, p.fail("workflow requires at least one connection")
	}
	close, err := p.expect(token.RBRACE)
	if err != nil {
		return nil, err
	}
	w.Span = Span{Start: start, End: close.offset + 1}
	w.Go, err = parseGo(p.source, w.Parameters, w.StateType, w.Prerequisite)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (p *dslParser) labelAhead() bool {
	if p.current().kind != token.IDENT {
		return false
	}
	i := p.index + 1
	for i < len(p.tokens) && p.tokens[i].kind == token.SEMICOLON && p.tokens[i].literal == "\n" {
		i++
	}
	return i < len(p.tokens) && p.tokens[i].kind == token.COLON
}

func (p *dslParser) directions(direction string) ([]Connection, error) {
	left, err := p.identifier()
	if err != nil {
		return nil, err
	}
	var connections []Connection
	for {
		first := p.current()
		if first.kind != token.DEC && first.kind != token.ARROW {
			return nil, p.fail("expected '-->' or '<--'")
		}
		p.advance()
		second := p.current()
		want := token.GTR
		if first.kind == token.ARROW {
			want = token.SUB
		}
		if second.kind != want || second.offset != first.offset+2 {
			return nil, p.fail("arrow must be contiguous: '-->' or '<--'")
		}
		p.advance()
		right, err := p.identifier()
		if err != nil {
			return nil, err
		}
		connection := Connection{From: left, To: right, Direction: direction, Span: Span{Start: left.Span.Start, End: right.Span.End}}
		if first.kind == token.ARROW {
			connection.From, connection.To = right, left
		}
		connections = append(connections, connection)
		left = right
		if p.current().kind != token.DEC && p.current().kind != token.ARROW {
			break
		}
	}
	return connections, nil
}

func (p *dslParser) branch() (Connection, error) {
	from, err := p.identifier()
	if err != nil {
		return Connection{}, err
	}
	first, err := p.expect(token.SUB)
	if err != nil {
		return Connection{}, err
	}
	label, err := p.expect(token.STRING)
	if err != nil {
		return Connection{}, err
	}
	second, err := p.expect(token.SUB)
	if err != nil {
		return Connection{}, err
	}
	last, err := p.expect(token.GTR)
	if err != nil {
		return Connection{}, err
	}
	// Read the original literal rather than using its normalized scanner value.
	// Unquote also rejects whitespace or comments between the string and '-'.
	direction, unquoteErr := strconv.Unquote(p.source.Text(Span{Start: label.offset, End: second.offset}))
	if label.offset != first.offset+1 || last.offset != second.offset+1 || unquoteErr != nil {
		return Connection{}, &scanner.Error{Pos: p.source.Position(first.offset), Msg: "branch arrow must be contiguous: -\"label\"-> or -`label`->"}
	}
	to, err := p.identifier()
	if err != nil {
		return Connection{}, err
	}
	return Connection{From: from, To: to, Direction: direction, Span: Span{Start: from.Span.Start, End: to.Span.End}}, nil
}
