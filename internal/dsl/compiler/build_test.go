package compiler

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// These tests invoke the real Go toolchain, including package name resolution.
// Parsing generated source alone cannot detect a value shadowing a Go type.
func TestGeneratedGoBuildAndRun(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", "..", ".."))
	require.NoError(t, err)
	dir := t.TempDir()
	writeBuildFile(t, dir, "go.mod", fmt.Sprintf("module example.org/dsltest\n\ngo 1.24.0\nrequire github.com/JSYoo5B/chain v0.0.0\nreplace github.com/JSYoo5B/chain => %q\n", root))
	writeBuildFile(t, dir, "model/model.go", "package model\ntype State struct { Value int }\n")
	const intHelper = `package sample
import ("context"; "github.com/JSYoo5B/chain")
func makeAction() chain.Action[int] {
    return chain.NewSimpleAction("increment", func(_ context.Context, v int) (int, error) { return v + 1, nil })
}`
	cases := []struct{ name, source, helper, call, input, want, imports string }{
		{"primitive", `package sample
workflow Build(int int) Flow[int] {
    prerequisite { step := makeAction() }
    nodes: step
    success: step --> end
}`, intHelper, "Build(0)", "41", "42", ""},
		{"local", `package sample
workflow Build() Flow[State] {
    prerequisite { State := 1; _ = State; step := makeAction() }
    nodes: step
    success: step --> end
}`, `package sample
import ("context"; "github.com/JSYoo5B/chain")
type State struct { Value int }
func makeAction() chain.Action[State] {
    return chain.NewSimpleAction("increment", func(_ context.Context, v State) (State, error) { v.Value++; return v, nil })
}`, "Build()", "State{Value: 41}", "State{Value: 42}", ""},
		{"qualified", `package sample
import item "example.org/dsltest/model"
workflow Build(item int) Flow[item.State] {
    prerequisite { step := makeAction() }
    nodes: step
    success: step --> end
}`, `package sample
import ("context"; "github.com/JSYoo5B/chain"; "example.org/dsltest/model")
func makeAction() chain.Action[model.State] {
    return chain.NewSimpleAction("increment", func(_ context.Context, v model.State) (model.State, error) { v.Value++; return v, nil })
}`, "Build(0)", "model.State{Value: 41}", "model.State{Value: 42}", `"example.org/dsltest/model"`},
		{"hygiene", `package sample
workflow Build(_chain, _workflow int) Flow[int] {
    prerequisite { _chain1 := 0; _workflow1 := 0; _, _ = _chain1, _workflow1; step := makeAction() }
    nodes: step
    success: step --> end
}`, intHelper, "Build(0, 0)", "41", "42", ""},
		{"label", `package sample
workflow Build() Flow[int] {
    prerequisite { step := makeAction(); goto done; done: }
    nodes: step
    success: step --> end
}`, intHelper, "Build()", "41", "42", ""},
		{"closure", `package sample
import ("context"; "github.com/JSYoo5B/chain")
workflow Build() Flow[int] {
    prerequisite {
        success := chain.NewSimpleAction("increment", func(_ context.Context, v int) (int, error) { return v + 1, nil })
    }
    nodes: success
    success: success --> end
}`, "package sample\n", "Build()", "41", "42", ""},
	}
	for _, tc := range cases {
		generated, err := Compile("flow.chain", []byte(tc.source))
		require.NoError(t, err, tc.name)
		writeBuildFile(t, dir, tc.name+"/flow.chain", tc.source)
		writeBuildFile(t, dir, tc.name+"/flow.go", string(generated))
		writeBuildFile(t, dir, tc.name+"/helper.go", tc.helper)
		writeBuildFile(t, dir, tc.name+"/flow_test.go", fmt.Sprintf(`package sample
import ("context"; "testing"; %s)
func TestRun(t *testing.T) {
    w := %s
    if err := w.ValidateGraph(); err != nil { t.Fatal(err) }
    got, err := w.Run(context.Background(), %s)
    if err != nil { t.Fatal(err) }
    if want := (%s); got != want { t.Fatalf("got %%v, want %%v", got, want) }
}`, tc.imports, tc.call, tc.input, tc.want))
	}
	cmd := exec.Command("go", "test", "-count=1", "./...")
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GOWORK=off", "GOFLAGS=-mod=mod")
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "%s", output)
}

func writeBuildFile(t *testing.T, dir, path, text string) {
	t.Helper()
	name := filepath.Join(dir, path)
	require.NoError(t, os.MkdirAll(filepath.Dir(name), 0755))
	require.NoError(t, os.WriteFile(name, []byte(text), 0644))
}

func TestGeneratedGoBuildDiagnostics(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", "..", ".."))
	require.NoError(t, err)
	dir := t.TempDir()
	writeBuildFile(t, dir, "go.mod", fmt.Sprintf("module example.org/dsltest\n\ngo 1.24.0\nrequire github.com/JSYoo5B/chain v0.0.0\nreplace github.com/JSYoo5B/chain => %q\n", root))
	cases := []struct{ name, source, marker string }{
		{"body", "package p\nworkflow Build() Flow[int] {\nprerequisite {\n step := missing()\n}\nnodes: step\nsuccess: step --> end\n}", "missing"},
		{"inline", "package p\nworkflow Build() Flow[int] {\nprerequisite { step := missing() }\nnodes: step\nsuccess: step --> end\n}", "missing"},
		{"directive", "package p\nworkflow Build() Flow[int] {\nprerequisite {\n//line spoof.go:900\n step := missing()\n}\nnodes: step\nsuccess: step --> end\n}", "missing"},
		{"parameter", "package p\nworkflow Build(a Missing) Flow[int] { nodes: a; success: a --> end }", "Missing"},
		{"state", "package p\nworkflow Build(a Action) Flow[Missing] { nodes: a; success: a --> end }", "Missing"},
		{"node", "package p\nworkflow Build() Flow[int] {\nnodes: missing\nsuccess: missing --> end\n}", "missing"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			generated, err := CompileWithOptions("flow.chain", []byte(tc.source), Options{SourceName: "../dsl/flow.chain"})
			require.NoError(t, err)
			writeBuildFile(t, dir, tc.name+"/flow.go", string(generated))
			cmd := exec.Command("go", "build", "./"+tc.name)
			cmd.Dir = dir
			cmd.Env = append(os.Environ(), "GOWORK=off", "GOFLAGS=-mod=mod")
			output, err := cmd.CombinedOutput()
			require.Error(t, err, "invalid Go must fail to build")
			require.Contains(t, string(output), "undefined: "+tc.marker)
			line := strings.Count(tc.source[:strings.Index(tc.source, tc.marker)], "\n") + 1
			require.Contains(t, string(output), fmt.Sprintf("flow.chain:%d:", line), "%s", output)
			require.NotContains(t, string(output), "flow.go:")
			require.NotContains(t, string(output), "spoof.go:")
		})
	}
}
