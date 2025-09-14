package parser

import (
	"fmt"
	"github.com/JSYoo5B/chain/internal/dsl/ast"
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
			importLine: strings.Join([]string{
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
			importLine: strings.Join([]string{
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
		"with single parameter constructors": {
			declares: []string{`workflow helloWorld(name string) HelloWorld[string]`},
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
			declares: []string{`workflow helloWorld(name string, age int) HelloWorld[string]`},
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
			declares: []string{`workflow helloWorld(first, last string) HelloWorld[string]`},
			expected: []ast.WorkflowDeclaration{
				{
					ConstructorName:   ast.CodeLocation{Text: "helloWorld"},
					ConstructorParams: ast.CodeLocation{Text: "first, last string"},
					WorkflowName:      ast.CodeLocation{Text: "HelloWorld"},
					WorkflowType:      ast.CodeLocation{Text: "string"},
				},
			},
		},
		"with multiple generic type": {
			declares: []string{`workflow helloWorld() HelloWorld[int8|int16]`},
			expected: []ast.WorkflowDeclaration{
				{
					ConstructorName:   ast.CodeLocation{Text: "helloWorld"},
					ConstructorParams: ast.CodeLocation{Text: ""},
					WorkflowName:      ast.CodeLocation{Text: "HelloWorld"},
					WorkflowType:      ast.CodeLocation{Text: "int8|int16"},
				},
			},
		},
		"multiple workflows": {
			declares: []string{
				`workflow helloWorld() HelloWorld[string]`,
				`workflow goodByeWorld() GoodByeWorld[string]`,
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

func TestAstBuilder_PrerequisiteBlock(t *testing.T) {
	type testCase struct {
		prerequisite string
	}

	testCases := map[string]testCase{
		"empty golang code": {
			prerequisite: ``,
		},
		"simple hello world": {
			prerequisite: strings.Join([]string{
				`hello := printAction("hello")`,
				`world := printAction("world")`,
			}, "\n"),
		},
		"package symbol": {
			prerequisite: `action := chain.NewSimpleAction()`,
		},
		"comment in front": {
			prerequisite: strings.Join([]string{
				`// comment in prerequisite`,
				`action := chain.NewSimpleAction()`,
			}, "\n"),
		},
		"comment in end": {
			prerequisite: strings.Join([]string{
				`action := chain.NewSimpleAction()`,
				`// comment in prerequisite`,
			}, "\n"),
		},
		"single-line closure": {
			prerequisite: `closure := func(_ context.Context, i int) (int, error) { return i, nil }`,
		},
		"multi-line closure": {
			prerequisite: strings.Join([]string{
				`closure := func(_ context.Context, i int) (int, error) {`,
				`    return i, nil`,
				`}`,
			}, "\n"),
		},
		"multi-line parameters": {
			prerequisite: strings.Join([]string{
				`a := chain.NewSimpleAction(`,
				`    "test",`,
				`    someFunc,`,
				`)`,
			}, "\n"),
		},
		"if statement": {
			prerequisite: strings.Join([]string{
				`if true {`,
				`    a := chain.NewSimpleAction()`,
				`}`,
			}, "\n"),
		},
		"struct literal": {
			prerequisite: strings.Join([]string{
				`a := printAction{`,
				`    member: Member`,
				`}`,
			}, "\n"),
		},
		"string contains curly": {
			prerequisite: `a := "{\"number\": 123}"`,
		},
		"string contains incomplete curly": {
			prerequisite: `a := "{\"number\": 123"`,
		},
		"string contains inline-comment": {
			prerequisite: `a := "//comment"`,
		},
		"string contains comment block": {
			prerequisite: `a := "/*comment*/"`,
		},
		"variable declaration": {
			prerequisite: `var action Action[int]`,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			input := strings.Join([]string{
				`package chain_test`,
				`import "github.com/JSYoo5B/chain"`,
				``,
				`workflow helloworld() HelloWorld[string] {`,
				`    prerequisite {`,
				tc.prerequisite,
				`    }`,
				`    nodes:`,
				`        a, b, c`,
				`}`,
			}, "\n")

			result := buildAstForTest(input)
			require.Equal(t, "chain_test", result.Package.Name)
			require.Len(t, result.Imports, 1)
			require.Equal(t, "github.com/JSYoo5B/chain", result.Imports[0].Path)
			require.Len(t, result.Workflows, 1)

			assert.Equal(t,
				strings.Trim(tc.prerequisite, " \t\n"),
				strings.Trim(result.Workflows[0].Prerequisite.Code, " \t\n"))
		})
	}
}

func TestAstBuilder_PrerequisiteBlock_multiple(t *testing.T) {
	type testCase struct {
		prerequisites []string
	}

	testCases := map[string]testCase{
		"two blocks": {
			prerequisites: []string{
				``,
				`action := chain.NewSimpleAction()`,
			},
		},
		"three blocks": {
			prerequisites: []string{
				`a := "{\"number\": 123}"`,
				`closure := func(_ context.Context, i int) (int, error) { return i, nil }`,
				`action := chain.NewSimpleAction()`,
			},
		},
		"same blocks": {
			prerequisites: []string{
				`a := "//comment"`,
				`a := "//comment"`,
				`a := "//comment"`,
				`a := "//comment"`,
			},
		},
		"some has multi-line": {
			prerequisites: []string{
				strings.Join([]string{
					`if true {`,
					`    a := chain.NewSimpleAction()`,
					`}`,
				}, "\n"),
				strings.Join([]string{
					`closure := func(_ context.Context, i int) (int, error) {`,
					`    return i, nil`,
					`}`,
				}, "\n"),
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			input := strings.Join([]string{
				`package chain_test`,
				`import "github.com/JSYoo5B/chain"`,
			}, "\n")

			for i, prerequisite := range tc.prerequisites {
				workflowDefine := strings.Join([]string{
					fmt.Sprintf("workflow wf%d() WF%d[string] {", i, i),
					"    prerequisite {",
					prerequisite,
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

			assert.Len(t, result.Workflows, len(tc.prerequisites))
			for i, workflow := range result.Workflows {
				assert.Equal(t, fmt.Sprintf("WF%d", i), workflow.WorkflowName.Text)
				assert.Equal(t,
					strings.Trim(tc.prerequisites[i], " \t\n"),
					strings.Trim(workflow.Prerequisite.Code, " \t\n"))
			}
		})
	}
}

func TestAstBuilder_NodesBlock(t *testing.T) {
	type testCase struct {
		nodesCode string
		nodeNames []string
	}

	testCases := map[string]testCase{
		"simple": {
			nodesCode: `a, b, c`,
			nodeNames: []string{`a`, `b`, `c`},
		},
		"multi-lines": {
			nodesCode: strings.Join([]string{
				`a, b, c,`,
				`d, e`,
			}, "\n"),
			nodeNames: []string{`a`, `b`, `c`, `d`, `e`},
		},
		"multi-line with comment": {
			nodesCode: strings.Join([]string{
				`a, b, c, // trail comment`,
				`d, e`,
			}, "\n"),
			nodeNames: []string{`a`, `b`, `c`, `d`, `e`},
		},
		"comment is multi-lined": {
			nodesCode: strings.Join([]string{
				`a, b,/*`,
				`d, e`,
				`*/ c`,
			}, "\n"),
			nodeNames: []string{`a`, `b`, `c`},
		},
		"variant node names": {
			nodesCode: `branch, even, odd1, odd2`,
			nodeNames: []string{`branch`, `even`, `odd1`, `odd2`},
		},
		"node contains comment": {
			nodesCode: `branch, even, odd1, /*comment*/ odd2`,
			nodeNames: []string{`branch`, `even`, `odd1`, `odd2`},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			input := strings.Join([]string{
				`package chain_test`,
				`import "github.com/JSYoo5B/chain"`,
				``,
				`workflow helloworld() HelloWorld[string] {`,
				`    prerequisite {`,
				`    }`,
				`    nodes:`,
				tc.nodesCode,
				`}`,
			}, "\n")

			result := buildAstForTest(input)
			require.Equal(t, "chain_test", result.Package.Name)
			require.Len(t, result.Imports, 1)
			require.Equal(t, "github.com/JSYoo5B/chain", result.Imports[0].Path)
			require.Len(t, result.Workflows, 1)

			assert.Len(t, result.Workflows[0].NodesBlock.Nodes, len(tc.nodeNames))
			for i, node := range result.Workflows[0].NodesBlock.Nodes {
				assert.Equal(t, tc.nodeNames[i], node.Name)
			}
		})
	}
}

func TestAstBuilder_Directions(t *testing.T) {
	type testCase struct {
		directionBlock string
		directions     []ast.DirectionStatement
	}

	testCases := map[string]testCase{
		"left to right": {
			directionBlock: `hello --> end`,
			directions: []ast.DirectionStatement{
				{FromNode: "hello", ToNode: "end"},
			},
		},
		"right to left": {
			directionBlock: `end <-- hello`,
			directions: []ast.DirectionStatement{
				{FromNode: "hello", ToNode: "end"},
			},
		},
		"no-spaces": {
			directionBlock: `hello-->world`,
			directions: []ast.DirectionStatement{
				{FromNode: "hello", ToNode: "world"},
			},
		},
		"left to right chain": {
			directionBlock: `hello --> world --> end`,
			directions: []ast.DirectionStatement{
				{FromNode: "hello", ToNode: "world"},
				{FromNode: "world", ToNode: "end"},
			},
		},
		"right to left chain": {
			directionBlock: `end <-- world <-- hello`,
			directions: []ast.DirectionStatement{
				{FromNode: "world", ToNode: "end"},
				{FromNode: "hello", ToNode: "world"},
			},
		},
		"direction combined": {
			directionBlock: `hello --> end <-- world`,
			directions: []ast.DirectionStatement{
				{FromNode: "hello", ToNode: "end"},
				{FromNode: "world", ToNode: "end"},
			},
		},
		"multi-lined direction": {
			directionBlock: strings.Join([]string{
				`hello --> world -->`,
				`end`,
			}, "\n"),
			directions: []ast.DirectionStatement{
				{FromNode: "hello", ToNode: "world"},
				{FromNode: "world", ToNode: "end"},
			},
		},
		"multi direction lines": {
			directionBlock: strings.Join([]string{
				`hello --> world`,
				`world --> end`,
			}, "\n"),
			directions: []ast.DirectionStatement{
				{FromNode: "hello", ToNode: "world"},
				{FromNode: "world", ToNode: "end"},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			input := strings.Join([]string{
				`package chain_test`,
				`import "github.com/JSYoo5B/chain"`,
				``,
				`workflow helloworld() HelloWorld[string] {`,
				`    prerequisite {`,
				`    }`,
				`    nodes:`,
				`        a, b, c`,
				`    success:`,
				tc.directionBlock,
				`    error:`,
				tc.directionBlock,
				`    abort:`,
				tc.directionBlock,
				`}`,
			}, "\n")

			result := buildAstForTest(input)
			require.Equal(t, "chain_test", result.Package.Name)
			require.Len(t, result.Imports, 1)
			require.Equal(t, "github.com/JSYoo5B/chain", result.Imports[0].Path)
			require.Len(t, result.Workflows, 1)
			require.Len(t, result.Workflows[0].NodesBlock.Nodes, 3)

			resultDirections := result.Workflows[0].Successes
			assert.Len(t, resultDirections, len(tc.directions))
			for i, direction := range resultDirections {
				assert.Equal(t, tc.directions[i].FromNode, direction.FromNode)
				assert.Equal(t, tc.directions[i].ToNode, direction.ToNode)
			}

			resultDirections = result.Workflows[0].Errors
			assert.Len(t, resultDirections, len(tc.directions))
			for i, direction := range resultDirections {
				assert.Equal(t, tc.directions[i].FromNode, direction.FromNode)
				assert.Equal(t, tc.directions[i].ToNode, direction.ToNode)
			}

			resultDirections = result.Workflows[0].Aborts
			assert.Len(t, resultDirections, len(tc.directions))
			for i, direction := range resultDirections {
				assert.Equal(t, tc.directions[i].FromNode, direction.FromNode)
				assert.Equal(t, tc.directions[i].ToNode, direction.ToNode)
			}
		})
	}
}
