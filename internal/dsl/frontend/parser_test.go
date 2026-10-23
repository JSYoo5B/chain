package frontend

import (
	"embed"
	"go/ast"
	"go/token"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/*.chain
var parserFixtures embed.FS

func connectionStrings(w *Workflow) []string {
	var result []string
	for _, connection := range w.Connections {
		result = append(result, connection.From.Name+" -"+connection.Direction+"-> "+connection.To.Name)
	}
	return result
}

func workflowSource(body string) string {
	return "package example\nworkflow build() Example[int] {\n" + body + "\n}\n"
}

func TestParseExistingFixtures(t *testing.T) {
	for _, name := range []string{"helloworld", "collatz"} {
		t.Run(name, func(t *testing.T) {
			text, err := parserFixtures.ReadFile("testdata/" + name + ".chain")
			require.NoError(t, err)
			file, err := Parse(name+".chain", text)
			require.NoError(t, err)
			for _, w := range file.Workflows {
				require.NotNil(t, w.Go)
				require.NotNil(t, w.Go.Prerequisite)
			}
			if name == "helloworld" {
				require.Len(t, file.Workflows, 1)
				assert.Equal(t, "hello", file.Header.Name.Name)
				assert.Equal(t, []string{"hello -success-> world", "world -success-> end"}, connectionStrings(file.Workflows[0]))
			} else {
				require.Len(t, file.Workflows, 2)
				assert.Equal(t, []string{
					"branch -even-> even", "branch -odd-> odd1", "even -success-> end",
					"odd1 -success-> odd2", "odd2 -success-> end",
				}, connectionStrings(file.Workflows[0]))
				assert.Equal(t, []string{
					"even -success-> end", "branch -even-> even", "branch -odd-> odd1",
					"odd1 -success-> odd2", "odd2 -success-> even",
				}, connectionStrings(file.Workflows[1]))
			}
		})
	}
}

func TestParseRequiredNodesAndConnections(t *testing.T) {
	cases := map[string]struct {
		body string
		err  string
	}{
		"single node":              {body: "nodes: decide\nbranches: decide -\"done\"-> end"},
		"single success":           {body: "nodes: a\nsuccess: a --> end"},
		"failure only":             {body: "nodes: a\nfailure: a --> end"},
		"empty extra sections":     {body: "nodes: a\nsuccess:\nfailure: a --> end\nabort:"},
		"missing nodes":            {body: "success: a --> end", err: "expected \"nodes\""},
		"empty nodes":              {body: "nodes:\nsuccess: a --> end", err: "nodes section requires at least one node"},
		"empty nodes before close": {body: "nodes:", err: "nodes section requires at least one node"},
		"no connection":            {body: "nodes: a", err: "at least one connection"},
		"empty sections":           {body: "nodes: a\nsuccess:\nbranches:", err: "at least one connection"},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			file, err := Parse("required.chain", []byte(workflowSource(tc.body)))
			if tc.err != "" {
				require.ErrorContains(t, err, tc.err)
				assert.Nil(t, file)
				return
			}
			require.NoError(t, err)
			require.Len(t, file.Workflows, 1)
			assert.Len(t, file.Workflows[0].Nodes, 1)
			assert.Len(t, file.Workflows[0].Connections, 1)
		})
	}
}

func TestParseDirectionalConnections(t *testing.T) {
	file, err := Parse("directions.chain", []byte(workflowSource(`nodes: a, b, c
success:
 a-->b<--c
failure:
 end <-- a
abort:
 c --> end
success:
 b --> end`)))
	require.NoError(t, err)
	w := file.Workflows[0]
	assert.Equal(t, []string{
		"a -success-> b", "c -success-> b", "a -failure-> end",
		"c -abort-> end", "b -success-> end",
	}, connectionStrings(w))
	assert.Equal(t, "b<--c", file.Source.Text(w.Connections[1].Span))
	assert.Equal(t, "c", file.Source.Text(w.Connections[1].From.Span))
}

func TestParseBranchStrings(t *testing.T) {
	cases := []struct{ arrow, direction string }{
		{`-"even"->`, "even"},
		{"-`odd`->", "odd"},
		{`-"a\"b\n\uAC00"->`, "a\"b\n가"},
		{"-`first\r\nsecond`->", "first\nsecond"},
	}
	for _, tc := range cases {
		t.Run(tc.arrow, func(t *testing.T) {
			file, err := Parse("branch.chain", []byte(workflowSource("nodes: decide\nbranches: decide"+tc.arrow+"end")))
			require.NoError(t, err)
			edge := file.Workflows[0].Connections[0]
			assert.Equal(t, "decide", edge.From.Name)
			assert.Equal(t, "end", edge.To.Name)
			assert.Equal(t, tc.direction, edge.Direction)
			assert.Equal(t, "decide"+tc.arrow+"end", file.Source.Text(edge.Span))
		})
	}
}

func TestParseGoAndContextualDSLWords(t *testing.T) {
	const text = "package example\n" +
		"workflow build(success, end int, callback func(x, y int) (int, error), options ...Option) Flow[map[string][]Box[int]] {\n" +
		"prerequisite {\n" +
		" nodes := `} success: -->`\n" +
		" branches := func() int { return success + end }\n" +
		" _ = nodes; _ = branches\n" +
		"}\n" +
		"nodes: success, branches\n" +
		"success: success --> branches --> end\n}"
	file, err := Parse("context.chain", []byte(text))
	require.NoError(t, err)
	w := file.Workflows[0]
	assert.Equal(t, "build", w.Constructor.Name)
	assert.Equal(t, "Flow", w.Name.Name)
	assert.Equal(t, "success", w.Go.Parameters.List[0].Names[0].Name)
	assert.Equal(t, "end", w.Go.Parameters.List[0].Names[1].Name)
	assert.IsType(t, &ast.Ellipsis{}, w.Go.Parameters.List[2].Type)
	assert.IsType(t, &ast.MapType{}, w.Go.StateType)
	assert.Equal(t, "map[string][]Box[int]", file.Source.Text(w.StateType.Content))
	assert.Equal(t, []string{"success -success-> branches", "branches -success-> end"}, connectionStrings(w))
}

func TestParseHeaderAndPhysicalLocations(t *testing.T) {
	const text = "//line elsewhere.go:900\npackage example\n" +
		"import (\n workflow \"example.org/renamed\"\n . `fmt`\n _ \"context\"\n)\n" +
		"import \"strings\"\n" +
		"workflow build() Flow[int] { nodes: a; success: a --> end }"
	file, err := Parse("header.chain", []byte(text))
	require.NoError(t, err)
	require.Len(t, file.Header.Imports, 4)
	assert.Equal(t, "workflow", file.Header.Imports[0].Name.Name)
	assert.Equal(t, ".", file.Header.Imports[1].Name.Name)
	assert.Equal(t, "_", file.Header.Imports[2].Name.Name)
	assert.Nil(t, file.Header.Imports[3].Name, "retain implicit imports without guessing package names")
	assert.Equal(t, file.Source.Position(strings.Index(text, "example\n")), file.HeaderPosition(file.Header.Name.Pos()))
	assert.Equal(t, file.Source.Position(strings.Index(text, "workflow \"")), file.HeaderPosition(file.Header.Imports[0].Name.Pos()))
	assert.Equal(t, token.Position{}, file.HeaderPosition(token.NoPos))
	assert.Equal(t, text[:strings.LastIndex(text, "workflow build")], file.Source.Text(file.HeaderSpan))
}

func TestParseSoftKeywordLineBreaks(t *testing.T) {
	file, err := Parse("lines.chain", []byte("package example\nworkflow\nbuild() Flow[int] {\nprerequisite\n{}\nnodes\n:\na\nsuccess\n:\na -->\nend\n}"))
	require.NoError(t, err)
	assert.Equal(t, []string{"a -success-> end"}, connectionStrings(file.Workflows[0]))
}

func TestParseRejectsInvalidFiles(t *testing.T) {
	cases := map[string]string{
		"missing package":                 "workflow build() Flow[int] { nodes: a; success: a --> end }",
		"invalid import":                  "package p\nimport (x 123)\n",
		"ordinary Go declaration":         "package p\nfunc helper() {}",
		"missing constructor parentheses": "package p\nworkflow build Flow[int] {}",
		"missing state brackets":          "package p\nworkflow build() Flow int {}",
		"missing workflow close":          "package p\nworkflow build() Flow[int] { nodes: a; success: a --> end",
		"unknown section":                 workflowSource("nodes: a\nother: a --> end"),
		"duplicate nodes section":         workflowSource("nodes: a\nnodes: a\nsuccess: a --> end"),
		"missing destination":             workflowSource("nodes: a\nsuccess: a -->"),
		"missing arrow":                   workflowSource("nodes: a\nsuccess: a end"),
		"spaced arrow":                    workflowSource("nodes: a\nsuccess: a -- > end"),
		"comment in arrow":                workflowSource("nodes: a\nsuccess: a --/*gap*/> end"),
		"spaced reverse arrow":            workflowSource("nodes: a\nsuccess: end <- - a"),
		"spaced branch prefix":            workflowSource("nodes: a\nbranches: a - \"done\"-> end"),
		"spaced branch suffix":            workflowSource("nodes: a\nbranches: a -\"done\" -> end"),
		"comment in branch":               workflowSource("nodes: a\nbranches: a -\"done\"/*gap*/-> end"),
		"spaced branch tip":               workflowSource("nodes: a\nbranches: a -\"done\"- > end"),
		"rune branch label":               workflowSource("nodes: a\nbranches: a -'x'-> end"),
		"unseparated connections":         workflowSource("nodes: a\nsuccess: a --> end a --> end"),
		"invalid Go parameter":            "package p\nworkflow build(a *) Flow[int] { nodes: a; success: a --> end }",
		"invalid Go state":                "package p\nworkflow build() Flow[int|string] { nodes: a; success: a --> end }",
		"invalid Go body":                 workflowSource("prerequisite { a := }\nnodes: a\nsuccess: a --> end"),
		"escaping return":                 workflowSource("prerequisite { if invalid { return } }\nnodes: a\nsuccess: a --> end"),
		"unterminated string":             workflowSource("prerequisite { a := \"unfinished }\nnodes: a\nsuccess: a --> end"),
	}
	for name, text := range cases {
		t.Run(name, func(t *testing.T) {
			file, err := Parse("invalid.chain", []byte(text))
			require.Error(t, err)
			assert.Nil(t, file)
			assert.Contains(t, err.Error(), "invalid.chain:")
		})
	}
}

func FuzzParse(f *testing.F) {
	for _, seed := range []string{
		"", "package p\n", workflowSource("nodes: a\nsuccess: a --> end"),
		workflowSource("nodes: decide\nbranches: decide -`done`-> end"),
		workflowSource("prerequisite { callback := func() { return }; _ = callback }\nnodes: a\nsuccess: a --> end"),
		"package p\nworkflow build() Flow[map[string] {",
	} {
		f.Add([]byte(seed))
	}
	f.Fuzz(func(t *testing.T, text []byte) {
		file, err := Parse("fuzz.chain", text)
		if err != nil {
			if file != nil {
				t.Fatal("parse error exposed a successful file")
			}
			return
		}
		for _, w := range file.Workflows {
			if len(w.Nodes) == 0 || len(w.Connections) == 0 || w.Go == nil {
				t.Fatal("successful workflow is missing nodes, connections, or Go syntax")
			}
		}
	})
}
