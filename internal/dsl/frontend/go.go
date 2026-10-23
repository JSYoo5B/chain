package frontend

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/scanner"
	"go/token"
)

// GoSyntax contains the Go AST for a constructor's embedded source. Parsing
// checks syntax and the DSL return rule, but does not resolve names or types.
type GoSyntax struct {
	Parameters   *ast.FieldList
	StateType    ast.Expr
	Prerequisite *ast.BlockStmt
	Comments     []*ast.CommentGroup

	file  *token.File
	input *goInput
}

// Position maps a Go AST position back to the physical DSL source location.
func (g *GoSyntax) Position(pos token.Pos) token.Position {
	if !pos.IsValid() {
		return token.Position{}
	}
	return g.input.position(g.file.Offset(pos))
}

func remapGoErrors(err error, position func(int) token.Position) error {
	var diagnostics scanner.ErrorList
	if list, ok := err.(scanner.ErrorList); ok {
		for _, diagnostic := range list {
			diagnostics.Add(position(diagnostic.Pos.Offset), diagnostic.Msg)
		}
	} else {
		diagnostics.Add(position(0), err.Error())
	}
	// go/parser sorts by positions adjusted by //line directives. Restore
	// physical source order after mapping all diagnostics back to the DSL.
	diagnostics.Sort()
	return diagnostics.Err()
}

// Each copied segment has an original span. Synthetic segments have an empty
// original span at the delimiter or boundary that should receive diagnostics.
type goSegment struct {
	generated Span
	original  Span
}

type goInput struct {
	source   *Source
	text     bytes.Buffer
	segments []goSegment
}

func (g *goInput) copy(span Span) {
	start := g.text.Len()
	g.text.Write(g.source.text[span.Start:span.End])
	g.segments = append(g.segments, goSegment{
		generated: Span{Start: start, End: g.text.Len()}, original: span,
	})
}

func (g *goInput) insert(text string, anchor int) {
	start := g.text.Len()
	g.text.WriteString(text)
	g.segments = append(g.segments, goSegment{
		generated: Span{Start: start, End: g.text.Len()},
		original:  Span{Start: anchor, End: anchor},
	})
}

func (g *goInput) position(offset int) token.Position {
	for _, segment := range g.segments {
		if offset >= segment.generated.Start && offset < segment.generated.End {
			original := segment.original.Start
			if segment.original.Start != segment.original.End {
				original += offset - segment.generated.Start
			}
			return g.source.Position(original)
		}
	}
	// Parse errors at synthetic EOF belong to the end of the last source region.
	last := g.segments[len(g.segments)-1]
	return g.source.Position(last.original.End)
}

// parseGo embeds the already delimited Go regions in one function declaration.
// Parameters retain their Go grouping, and the bracketed state type becomes a
// single unnamed Go result. The wrapper is only for parsing, not execution.
func parseGo(source *Source, parameters, stateType GoRegion, prerequisite *GoRegion) (*GoSyntax, error) {
	input := &goInput{source: source}
	input.insert("package __chain\nfunc __chain", parameters.Span.Start)
	input.copy(parameters.Span)
	input.insert(" (", stateType.Span.Start)
	input.copy(stateType.Content)
	input.insert(")", stateType.Content.End)
	if prerequisite != nil {
		input.insert(" ", prerequisite.Span.Start)
		input.copy(prerequisite.Span)
		input.insert("\n", prerequisite.Span.End)
	} else {
		input.insert(" {}\n", stateType.Span.End)
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "__chain_embedded.go", input.text.Bytes(),
		parser.AllErrors|parser.ParseComments|parser.SkipObjectResolution)
	if err != nil {
		return nil, remapGoErrors(err, input.position)
	}

	fn := file.Decls[0].(*ast.FuncDecl)
	var diagnostics scanner.ErrorList
	if len(fn.Type.Results.List) != 1 || len(fn.Type.Results.List[0].Names) != 0 {
		diagnostics.Add(source.Position(stateType.Content.Start), "expected one unnamed Go type for workflow state")
		return nil, diagnostics.Err()
	}

	result := &GoSyntax{
		Parameters: fn.Type.Params,
		StateType:  fn.Type.Results.List[0].Type,
		Comments:   file.Comments,
		file:       fset.File(fn.Pos()),
		input:      input,
	}
	if prerequisite != nil {
		result.Prerequisite = fn.Body
		ast.Inspect(fn.Body, func(node ast.Node) bool {
			switch node := node.(type) {
			case *ast.FuncLit:
				// A return here exits a separate function, not the prerequisite.
				return false
			case *ast.ReturnStmt:
				diagnostics.Add(result.Position(node.Return), "return cannot exit a prerequisite block")
			}
			return true
		})
	}
	if err := diagnostics.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
