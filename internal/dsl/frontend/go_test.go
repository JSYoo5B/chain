package frontend

import (
	"errors"
	"go/ast"
	"go/scanner"
	"go/token"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type goFixture struct {
	source       *Source
	parameters   GoRegion
	stateType    GoRegion
	prerequisite *GoRegion
}

func newGoFixture(t *testing.T, parameters, stateType, body string) goFixture {
	t.Helper()
	text := "package example\n\nworkflow build" + parameters + " Example[" + stateType + "] {\n"
	if body != "" {
		text += "prerequisite " + body + "\n"
	}
	text += "nodes:\n a, b\nsuccess:\n a --> b --> end\n}\n"
	f := goFixture{source: NewSource("example.chain", []byte(text))}
	c, err := scan(f.source)
	require.NoError(t, err)
	for c.current().kind != token.LPAREN && c.current().kind != token.EOF {
		c.advance()
	}
	f.parameters, err = c.takeGoRegion()
	require.NoError(t, err)
	require.Equal(t, "Example", c.current().literal)
	c.advance()
	f.stateType, err = c.takeGoRegion()
	require.NoError(t, err)
	require.Equal(t, token.LBRACE, c.current().kind)
	c.advance()
	if body != "" {
		require.Equal(t, "prerequisite", c.current().literal)
		c.advance()
		region, err := c.takeGoRegion()
		require.NoError(t, err)
		f.prerequisite = &region
	}
	return f
}

func (f goFixture) parse() (*GoSyntax, error) {
	return parseGo(f.source, f.parameters, f.stateType, f.prerequisite)
}

func (f goFixture) positionOf(text string) token.Position {
	return f.source.Position(strings.Index(string(f.source.text), text))
}

func TestParseGoConstructor(t *testing.T) {
	fixture := newGoFixture(t,
		"(a, b int, callback func(x, y int) (int, error), actions map[string]chain.Action[State], options ...Option)",
		"map[string][]pkg.Box[int]",
		"{ success := callback; end := success; _ = end }",
	)
	syntax, err := fixture.parse()
	require.NoError(t, err)
	require.Len(t, syntax.Parameters.List, 4)
	assert.Equal(t, "a", syntax.Parameters.List[0].Names[0].Name)
	assert.Equal(t, "b", syntax.Parameters.List[0].Names[1].Name)
	assert.IsType(t, &ast.FuncType{}, syntax.Parameters.List[1].Type)
	assert.IsType(t, &ast.MapType{}, syntax.Parameters.List[2].Type)
	assert.IsType(t, &ast.Ellipsis{}, syntax.Parameters.List[3].Type)
	assert.IsType(t, &ast.MapType{}, syntax.StateType)
	require.NotNil(t, syntax.Prerequisite)
	assert.Len(t, syntax.Prerequisite.List, 3)
	assert.Equal(t, fixture.positionOf("a, b"), syntax.Position(syntax.Parameters.List[0].Pos()))
	assert.Equal(t, fixture.positionOf("map[string][]"), syntax.Position(syntax.StateType.Pos()))
	assert.Equal(t, fixture.source.Position(fixture.prerequisite.Span.Start), syntax.Position(syntax.Prerequisite.Lbrace))
}

func TestParseGoWithoutPrerequisite(t *testing.T) {
	fixture := newGoFixture(t, "()", "State", "")
	syntax, err := fixture.parse()
	require.NoError(t, err)
	assert.Empty(t, syntax.Parameters.List)
	assert.Nil(t, syntax.Prerequisite)
	assert.Equal(t, "State", syntax.StateType.(*ast.Ident).Name)
	assert.Equal(t, token.Position{}, syntax.Position(token.NoPos))
}

func TestPrerequisiteAllowsFunctionReturns(t *testing.T) {
	bodies := map[string]string{
		"empty":                "{}",
		"closure":              "{ callback := func() int { return 1 }; _ = callback }",
		"nested closures":      "{ callback := func() func() int { return func() int { return 1 } }; _ = callback }",
		"immediate invocation": "{ value := func() int { return 1 }(); _ = value }",
		"deferred function":    "{ defer func() { return }() }",
		"goroutine":            "{ go func() { return }() }",
		"closure in condition": "{ if func() bool { return true }() { work() } }",
		"named result":         "{ callback := func() (value int) { return }; _ = callback }",
		"text and comments":    "{ text := `return`; /* return */ _ = text }",
	}
	for name, body := range bodies {
		t.Run(name, func(t *testing.T) {
			syntax, err := newGoFixture(t, "()", "State", body).parse()
			require.NoError(t, err)
			require.NotNil(t, syntax.Prerequisite)
		})
	}
}

func TestPrerequisiteRejectsEscapingReturn(t *testing.T) {
	bodies := map[string]string{
		"direct":     "{ return }",
		"with value": "{ return value }",
		"block":      "{ { return } }",
		"if":         "{ if invalid { return } }",
		"for":        "{ for { return } }",
		"switch":     "{ switch { case invalid: return } }",
		"select":     "{ select { default: return } }",
		"label":      "{ done: return }",
	}
	for name, body := range bodies {
		t.Run(name, func(t *testing.T) {
			fixture := newGoFixture(t, "()", "State", body)
			syntax, err := fixture.parse()
			require.Error(t, err)
			assert.Nil(t, syntax)
			var diagnostics scanner.ErrorList
			require.True(t, errors.As(err, &diagnostics))
			require.Len(t, diagnostics, 1)
			assert.Equal(t, fixture.positionOf("return"), diagnostics[0].Pos)
			assert.Equal(t, "return cannot exit a prerequisite block", diagnostics[0].Msg)
		})
	}
}

func TestPrerequisiteReportsOnlyOuterReturns(t *testing.T) {
	body := "{\n" +
		"callback := func() int { if invalid { return 1 }; return 2 }\n" +
		"if invalid { return /* first */ }\n" +
		"for { return /* second */ }\n" +
		"}"
	fixture := newGoFixture(t, "()", "State", body)
	syntax, err := fixture.parse()
	assert.Nil(t, syntax)
	var diagnostics scanner.ErrorList
	require.True(t, errors.As(err, &diagnostics))
	require.Len(t, diagnostics, 2)
	assert.Equal(t, fixture.positionOf("return /* first */"), diagnostics[0].Pos)
	assert.Equal(t, fixture.positionOf("return /* second */"), diagnostics[1].Pos)
}

func TestGoASTPositionsPreserveOriginalBytes(t *testing.T) {
	fixture := newGoFixture(t,
		"(\r\n // parameter comment\r\n 값 int,\r\n callback func(int) error,\r\n)",
		"map[string][]pkg.Box[int]",
		"{\r\n//line elsewhere.go:900\r\n success := func() int { return 1 }\r\n _ = success\r\n}",
	)
	syntax, err := fixture.parse()
	require.NoError(t, err)
	assert.Equal(t, fixture.positionOf("값"), syntax.Position(syntax.Parameters.List[0].Names[0].Pos()))
	assert.Equal(t, fixture.positionOf("map[string][]"), syntax.Position(syntax.StateType.Pos()))
	require.Len(t, syntax.Comments, 2)
	assert.Equal(t, fixture.positionOf("// parameter comment"), syntax.Position(syntax.Comments[0].Pos()))
	ast.Inspect(syntax.Prerequisite, func(node ast.Node) bool {
		if node, ok := node.(*ast.ReturnStmt); ok {
			assert.Equal(t, fixture.positionOf("return 1"), syntax.Position(node.Return))
		}
		return true
	})
	assert.Contains(t, fixture.source.Text(fixture.prerequisite.Content), "\r\n")
}

func TestGoSyntaxErrorsUseOriginalLocations(t *testing.T) {
	cases := []struct {
		name, parameters, stateType, body, errorAt string
	}{
		{"parameter", "(a *)", "State", "{}", ") Example"},
		{"state", "()", "int8|int16", "{}", "|int16"},
		{"body", "()", "State", "{\n value :=\n}", "}\nnodes:"},
		{"line directive", "()", "State", "{\n//line elsewhere.go:900\n value :=\n}", "}\nnodes:"},
		{"synthetic type delimiter", "()", "map[string]", "{}", "] {"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fixture := newGoFixture(t, tc.parameters, tc.stateType, tc.body)
			syntax, err := fixture.parse()
			assert.Nil(t, syntax)
			var diagnostics scanner.ErrorList
			require.True(t, errors.As(err, &diagnostics))
			require.NotEmpty(t, diagnostics)
			assert.Equal(t, fixture.positionOf(tc.errorAt), diagnostics[0].Pos)
			for _, diagnostic := range diagnostics {
				assert.Equal(t, "example.chain", diagnostic.Pos.Filename)
			}
		})
	}
}

func TestWorkflowStateRequiresOneUnnamedType(t *testing.T) {
	for _, stateType := range []string{"", "int, string", "value int", "a, b int"} {
		t.Run(stateType, func(t *testing.T) {
			fixture := newGoFixture(t, "()", stateType, "{}")
			syntax, err := fixture.parse()
			assert.Nil(t, syntax)
			require.ErrorContains(t, err, "expected one unnamed Go type for workflow state")
		})
	}
}
