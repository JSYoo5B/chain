package parser

import (
	"fmt"
	"github.com/JSYoo5B/chain/internal/compiler/ast"
	"github.com/antlr4-go/antlr/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func buildAstForTest(input string) *ast.ChainCode {
	is := antlr.NewInputStream(input)
	lexer := NewCommonLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parser := NewChainParser(stream)

	parser.RemoveErrorListeners()

	tree := parser.SourceFile()

	builder := NewAstBuilder(stream)
	antlr.ParseTreeWalkerDefault.Walk(builder, tree)
	return &builder.Result
}

func TestAstBuilder_Package(t *testing.T) {
	input := `package chain_test`

	result := buildAstForTest(input)

	assert.Equal(t, "chain_test", result.Package.Name)
	fmt.Println(result)
}

func TestAstBuilder_Import(t *testing.T) {
	type testCase struct {
		importLine string
		expected   []ast.Import
	}

	testCases := map[string]testCase{
		"no import": {
			importLine: "",
			expected:   nil,
		},
		"single import": {
			importLine: `import "context"`,
			expected: []ast.Import{
				{Alias: "context", Path: "context"},
			},
		},
		"single import with backquote": {
			importLine: "import `fmt`",
			expected: []ast.Import{
				{Alias: "fmt", Path: "fmt"},
			},
		},
		"single import with dot import": {
			importLine: `import . "fmt"`,
			expected: []ast.Import{
				{Alias: ".", Path: "fmt"},
			},
		},
		"single import with alias": {
			importLine: `import m "math"`,
			expected: []ast.Import{
				{Alias: "m", Path: "math"},
			},
		},
		"single import with ignore alias": {
			importLine: `import _ "ignore"`,
			expected: []ast.Import{
				{Alias: "_", Path: "ignore"},
			},
		},
		"single import with remote package": {
			importLine: `import "github.com/JSYoo5B/chain"`,
			expected: []ast.Import{
				{Alias: "chain", Path: "github.com/JSYoo5B/chain"},
			},
		},
		"single import with remote package version path": {
			importLine: `import "github.com/antlr4-go/antlr/v4"`,
			expected: []ast.Import{
				{Alias: "antlr", Path: "github.com/antlr4-go/antlr/v4"},
			},
		},
		"single import with remote package version dot": {
			importLine: `import "gopkg.in/yaml.v3"`,
			expected: []ast.Import{
				{Alias: "yaml", Path: "gopkg.in/yaml.v3"},
			},
		},
		"blocked multi import": {
			importLine: strings.Join(
				[]string{
					`import (`,
					`    "context"`,
					`    "fmt"`,
					`)`,
				}, "\n"),
			expected: []ast.Import{
				{Alias: "context", Path: "context"},
				{Alias: "fmt", Path: "fmt"},
			},
		},
		"unblocked multi import": {
			importLine: strings.Join(
				[]string{
					`import "context"`,
					`import "fmt"`,
				}, "\n"),
			expected: []ast.Import{
				{Alias: "context", Path: "context"},
				{Alias: "fmt", Path: "fmt"},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			input := fmt.Sprintf("package chain_test\n%s", tc.importLine)

			result := buildAstForTest(input)
			require.Equal(t, "chain_test", result.Package.Name)

			assert.Len(t, result.Imports, len(tc.expected))
			for i, imp := range result.Imports {
				assert.Equal(t, tc.expected[i].Alias, imp.Alias)
				assert.Equal(t, tc.expected[i].Path, imp.Path)
			}
		})
	}
}

func TestAstBuilder_WorkflowDef(t *testing.T) {
	type testCase struct {
		declares []string
		expected []ast.WorkflowDeclaration
	}

	testCases := map[string]testCase{
		"simple constructor": {
			declares: []string{`workflow helloWorld() generates HelloWorld[string]`},
			expected: []ast.WorkflowDeclaration{
				{
					ConstructorName:   ast.CodeLocation{Text: "helloWorld"},
					ConstructorParams: ast.CodeLocation{Text: ""},
					WorkflowName:      ast.CodeLocation{Text: "HelloWorld"},
					WorkflowType:      ast.CodeLocation{Text: "string"},
				},
			},
		},
		"with single parameter constructors": {
			declares: []string{`workflow helloWorld(name string) generates HelloWorld[string]`},
			expected: []ast.WorkflowDeclaration{
				{
					ConstructorName:   ast.CodeLocation{Text: "helloWorld"},
					ConstructorParams: ast.CodeLocation{Text: "name string"},
					WorkflowName:      ast.CodeLocation{Text: "HelloWorld"},
					WorkflowType:      ast.CodeLocation{Text: "string"},
				},
			},
		},
		"with multiple parameter constructors": {
			declares: []string{`workflow helloWorld(name string, age int) generates HelloWorld[string]`},
			expected: []ast.WorkflowDeclaration{
				{
					ConstructorName:   ast.CodeLocation{Text: "helloWorld"},
					ConstructorParams: ast.CodeLocation{Text: "name string, age int"},
					WorkflowName:      ast.CodeLocation{Text: "HelloWorld"},
					WorkflowType:      ast.CodeLocation{Text: "string"},
				},
			},
		},
		"with parameter signature type reuse": {
			declares: []string{`workflow helloWorld(firstName, lastName string) generates HelloWorld[string]`},
			expected: []ast.WorkflowDeclaration{
				{
					ConstructorName:   ast.CodeLocation{Text: "helloWorld"},
					ConstructorParams: ast.CodeLocation{Text: "firstName, lastName string"},
					WorkflowName:      ast.CodeLocation{Text: "HelloWorld"},
					WorkflowType:      ast.CodeLocation{Text: "string"},
				},
			},
		},
		"with multiple generic type": {
			declares: []string{`workflow helloWorld() generates HelloWorld[int8|int16]`},
			expected: []ast.WorkflowDeclaration{
				{
					ConstructorName:   ast.CodeLocation{Text: "helloWorld"},
					ConstructorParams: ast.CodeLocation{Text: ""},
					WorkflowName:      ast.CodeLocation{Text: "HelloWorld"},
					WorkflowType:      ast.CodeLocation{Text: "int8|int16"},
				},
			},
		},
		"skip generates": {
			declares: []string{`workflow helloWorld() HelloWorld[string]`},
			expected: []ast.WorkflowDeclaration{
				{
					ConstructorName:   ast.CodeLocation{Text: "helloWorld"},
					ConstructorParams: ast.CodeLocation{Text: ""},
					WorkflowName:      ast.CodeLocation{Text: "HelloWorld"},
					WorkflowType:      ast.CodeLocation{Text: "string"},
				},
			},
		},
		"multiple workflows": {
			declares: []string{
				`workflow helloWorld() generates HelloWorld[string]`,
				`workflow goodByeWorld() generates GoodByeWorld[string]`,
			},
			expected: []ast.WorkflowDeclaration{
				{
					ConstructorName:   ast.CodeLocation{Text: "helloWorld"},
					ConstructorParams: ast.CodeLocation{Text: ""},
					WorkflowName:      ast.CodeLocation{Text: "HelloWorld"},
					WorkflowType:      ast.CodeLocation{Text: "string"},
				},
				{
					ConstructorName:   ast.CodeLocation{Text: "goodByeWorld"},
					ConstructorParams: ast.CodeLocation{Text: ""},
					WorkflowName:      ast.CodeLocation{Text: "GoodByeWorld"},
					WorkflowType:      ast.CodeLocation{Text: "string"},
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			input := strings.Join([]string{
				`package chain_test`,
				`import "github.com/JSYoo5B/chain"`,
			}, "\n")

			for _, declare := range tc.declares {
				workflowDefine := strings.Join([]string{
					fmt.Sprintf("%s {", declare),
					"    prerequisite {",
					"    }",
					"    nodes:",
					"        a, b, c",
					"}",
				}, "\n")
				input += workflowDefine + "\n"
			}

			result := buildAstForTest(input)
			require.Equal(t, "chain_test", result.Package.Name)
			require.Len(t, result.Imports, 1)
			require.Equal(t, "github.com/JSYoo5B/chain", result.Imports[0].Path)

			assert.Len(t, result.Workflows, len(tc.expected))
			for i, workflow := range result.Workflows {
				assert.Equal(t, tc.expected[i].ConstructorName.Text, workflow.ConstructorName.Text)
				assert.Equal(t, tc.expected[i].ConstructorParams.Text, workflow.ConstructorParams.Text)
				assert.Equal(t, tc.expected[i].WorkflowName.Text, workflow.WorkflowName.Text)
				assert.Equal(t, tc.expected[i].WorkflowType.Text, workflow.WorkflowType.Text)
			}
		})
	}
}
