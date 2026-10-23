package compiler

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGeneratedCollatzIsCurrent(t *testing.T) {
	dir := filepath.Join("..", "..", "..", "examples", "branch")
	source, err := os.ReadFile(filepath.Join(dir, "collatz.chain"))
	require.NoError(t, err)
	generated, err := Compile("collatz.chain", source)
	require.NoError(t, err)
	checkedIn, err := os.ReadFile(filepath.Join(dir, "collatz_chain.gen.go"))
	require.NoError(t, err)
	assert.Equal(t, string(checkedIn), string(generated), "run go generate in examples/branch")
	again, err := Compile("collatz.chain", source)
	require.NoError(t, err)
	assert.Equal(t, generated, again, "generation must be deterministic")
}

func TestCompileAvoidsGeneratedNameCollisions(t *testing.T) {
	const source = `package scope
import _chain "github.com/JSYoo5B/chain"
workflow build(_workflow _chain.Action[int], _workflow1 int) Scope[int] {
    prerequisite { _chain1 := 0; _ = _chain1 }
    nodes: _workflow
    success: _workflow --> end
}`
	generated, err := Compile("scope.chain", []byte(source))
	require.NoError(t, err)
	_, err = parser.ParseFile(token.NewFileSet(), "scope.go", generated, parser.AllErrors)
	require.NoError(t, err)
	assert.Contains(t, string(generated), `import _chain2 "github.com/JSYoo5B/chain"`)
	assert.Contains(t, string(generated), "_workflow2 := _chain2.NewWorkflow[int]")
	assert.Contains(t, string(generated), `_chain2.Terminate[int]()`)
}

func TestCompileQuotesBranchDirections(t *testing.T) {
	const source = "package quoted\nworkflow build(a Action) Quoted[int] {\nnodes: a\nbranches: a -`a\"b\nline`-> end\n}"
	generated, err := Compile("quoted.chain", []byte(source))
	require.NoError(t, err)
	_, err = parser.ParseFile(token.NewFileSet(), "quoted.go", generated, parser.AllErrors)
	require.NoError(t, err)
	assert.Contains(t, string(generated), `"a\"b\nline":`)
}

func TestCompilePreservesTrailingLabel(t *testing.T) {
	const source = `package labeled
workflow build(a Action) Labeled[int] {
    prerequisite { goto done; done: }
    nodes: a
    success: a --> end
}`
	generated, err := Compile("labeled.chain", []byte(source))
	require.NoError(t, err)
	_, err = parser.ParseFile(token.NewFileSet(), "labeled.go", generated, parser.AllErrors)
	require.NoError(t, err)
	assert.Contains(t, string(generated), "done:\n\t;")
}

func TestCompileDoesNotReturnCodeOnError(t *testing.T) {
	for _, source := range []string{
		"package p\nworkflow broken() Flow[int] {}",
		"package p\nworkflow broken() Flow[int] { nodes: a; success: a --> missing }",
		"package p\nworkflow broken() Flow[int] { prerequisite { return }; nodes: a; success: a --> end }",
	} {
		generated, err := Compile("broken.chain", []byte(source))
		require.Error(t, err)
		assert.Nil(t, generated)
	}
}
