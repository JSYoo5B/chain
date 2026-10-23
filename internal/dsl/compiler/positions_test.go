package compiler

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGeneratedSourcePositions(t *testing.T) {
	for name, body := range map[string]string{
		"multiline":         "{\n step := missing()\n}",
		"inline statements": "{ a := 1; _ = a; step := missing() }",
		"closure":           "{ step := factory(func() int { return missing() }) }",
		"line directive":    "{\n//line spoof.go:900\n step := missing()\n}",
		"block directive":   "{ step := /*line spoof.go:900*/ missing() }",
		"CRLF":              "{\r\n _ = `text\r\n} nodes:`\r\n step := missing()\r\n}",
	} {
		t.Run(name, func(t *testing.T) {
			source := "package p\nworkflow Build() Flow[int] {\nprerequisite " + body + "\nnodes: step\nsuccess: step --> end\n}"
			generated, err := Compile("/original/flow.chain", []byte(source))
			require.NoError(t, err)
			positions := token.NewFileSet()
			file, err := parser.ParseFile(positions, "generated.go", generated, parser.SkipObjectResolution)
			require.NoError(t, err)
			found := false
			ast.Inspect(file, func(node ast.Node) bool {
				if ident, ok := node.(*ast.Ident); ok && ident.Name == "missing" {
					pos := positions.Position(ident.Pos())
					assert.Equal(t, "flow.chain", pos.Filename)
					assert.Equal(t, strings.Count(source[:strings.Index(source, "missing")], "\n")+1, pos.Line)
					found = true
				}
				return true
			})
			require.True(t, found)
		})
	}
}

func TestGeneratedSourcePath(t *testing.T) {
	source := []byte("package p\nworkflow Build(a Action) Flow[int] { nodes: a; success: a --> end }")
	generated, err := CompileWithOptions("flow.chain", source, Options{SourceName: "../dsl/flow.chain"})
	require.NoError(t, err)
	assert.Contains(t, string(generated), "//line ../dsl/flow.chain:")
	_, err = CompileWithOptions("flow.chain", source, Options{SourceName: "bad\nname.chain"})
	require.ErrorContains(t, err, "filename contains")
}
