package compiler

import (
	"errors"
	"go/scanner"
	"strings"
	"testing"

	"github.com/JSYoo5B/chain/internal/dsl/frontend"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func parseWorkflow(t *testing.T, body string) *frontend.File {
	t.Helper()
	file, err := frontend.Parse("graph.chain", []byte("package example\nworkflow build() Flow[int] {\n"+body+"\n}\n"))
	require.NoError(t, err)
	return file
}

func TestAnalyzeValidGraphs(t *testing.T) {
	for name, body := range map[string]string{
		"single node":         "nodes: a\nsuccess: a --> end",
		"branch termination":  "nodes: a\nbranches: a -\"done\"-> end",
		"diamond":             "nodes: a, b, c, d\nbranches: a -\"b\"-> b\na -\"c\"-> c\nsuccess: b --> d\nc --> d --> end",
		"reverse connection":  "nodes: a, b\nsuccess: a <-- b",
		"all direction types": "nodes: a, b, c\nsuccess: a --> b\nfailure: a --> c\nabort: b --> c",
	} {
		t.Run(name, func(t *testing.T) {
			file := parseWorkflow(t, body)
			program, err := Analyze(file)
			require.NoError(t, err)
			require.Len(t, program.Workflows, 1)
			assert.Equal(t, "a", program.Workflows[0].Nodes[0].Name.Name)
			assert.Len(t, program.Workflows[0].Nodes, len(file.Workflows[0].Nodes))
		})
	}
}

func TestAnalyzeRejectsInvalidGraphs(t *testing.T) {
	cases := []struct{ name, body, message string }{
		{"duplicate node", "nodes: a, a\nsuccess: a --> end", "duplicate node"},
		{"end member", "nodes: a, end\nsuccess: a --> end", "termination target"},
		{"blank member", "nodes: _\nsuccess: _ --> end", "blank identifier"},
		{"end source", "nodes: a\nsuccess: end --> a", "cannot be the source"},
		{"unknown source", "nodes: a\nsuccess: missing --> a", "undeclared source"},
		{"unknown target", "nodes: a\nsuccess: a --> missing", "undeclared destination"},
		{"duplicate direction", "nodes: a, b, c\nsuccess: a --> b\na --> c", "duplicate direction"},
		{"repeated same edge", "nodes: a\nsuccess: a --> end\na --> end", "duplicate direction"},
		{"repeated section", "nodes: a, b\nsuccess: a --> b\nsuccess: a --> end", "duplicate direction"},
		{"builtin branch collision", "nodes: a\nsuccess: a --> end\nbranches: a -\"success\"-> end", "duplicate direction"},
		{"duplicate branch string", "nodes: a\nbranches: a -\"done\"-> end\na -`done`-> end", "duplicate direction"},
		{"self loop", "nodes: a\nsuccess: a --> a", "directed cycle"},
		{"cross direction cycle", "nodes: a, b\nsuccess: a --> b\nfailure: b --> a", "directed cycle"},
		{"disconnected cycle", "nodes: a, b, c\nsuccess: a --> end\nb --> c --> b", "directed cycle"},
		{"disconnected terminal nodes", "nodes: a, b\nsuccess: a --> end\nb --> end", "disconnected"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			program, err := Analyze(parseWorkflow(t, tc.body))
			require.ErrorContains(t, err, tc.message)
			assert.Nil(t, program)
			assert.Contains(t, err.Error(), "graph.chain:")
		})
	}
}

func TestAnalyzeDiagnosticsUseOriginalLocations(t *testing.T) {
	file := parseWorkflow(t, "nodes: a\nsuccess: a --> missing")
	_, err := Analyze(file)
	var diagnostics scanner.ErrorList
	require.True(t, errors.As(err, &diagnostics))
	require.Len(t, diagnostics, 1)
	assert.Equal(t, file.Source.Position(strings.Index(file.Source.Text(file.Workflows[0].Span), "missing")+file.Workflows[0].Span.Start), diagnostics[0].Pos)
}
