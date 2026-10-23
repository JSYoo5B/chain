package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunWritesOnlySuccessfulCompilation(t *testing.T) {
	dir := t.TempDir()
	input, output := filepath.Join(dir, "flow.chain"), filepath.Join(dir, "flow.go")
	require.NoError(t, os.WriteFile(input, []byte("package p\nworkflow build(a Action) Flow[int] { nodes: a; success: a --> end }"), 0644))
	var stdout, stderr bytes.Buffer
	require.NoError(t, run([]string{input}, &stdout, &stderr))
	assert.Contains(t, stdout.String(), "func build(a Action)")
	require.NoError(t, run([]string{"-o", output, input}, &stdout, &stderr))
	generated, err := os.ReadFile(output)
	require.NoError(t, err)
	assert.Equal(t, stdout.Bytes(), generated)
	// A distinct existing output must remain writable on subsequent runs.
	require.NoError(t, run([]string{"-o", output, input}, &stdout, &stderr))

	require.NoError(t, os.WriteFile(input, []byte("package p\nworkflow broken() Flow[int] {}"), 0644))
	stdout.Reset()
	require.Error(t, run([]string{input}, &stdout, &stderr))
	assert.Empty(t, stdout.Bytes())
	require.Error(t, run([]string{"-o", output, input}, &stdout, &stderr))
	afterFailure, err := os.ReadFile(output)
	require.NoError(t, err)
	assert.Equal(t, generated, afterFailure)
}

func TestRunRejectsInvalidArguments(t *testing.T) {
	var stdout, stderr bytes.Buffer
	require.NoError(t, run([]string{"-h"}, &stdout, &stderr))
	require.ErrorContains(t, run(nil, &stdout, &stderr), "usage:")
	require.Error(t, run([]string{"one.chain", "two.chain"}, &stdout, &stderr))
	require.Error(t, run([]string{"-unknown"}, &stdout, &stderr))
	require.ErrorContains(t, run([]string{"-o", "same.chain", "same.chain"}, &stdout, &stderr), "must differ")
	assert.Empty(t, stdout.Bytes())
}

func TestRunMapsSourceRelativeToOutput(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "dsl"), 0755))
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "generated"), 0755))
	input := filepath.Join(dir, "dsl", "flow.chain")
	output := filepath.Join(dir, "generated", "flow.go")
	require.NoError(t, os.WriteFile(input, []byte("package p\nworkflow Build(a Action) Flow[int] { nodes: a; success: a --> end }"), 0644))
	var stdout, stderr bytes.Buffer
	require.NoError(t, run([]string{"-o", output, input}, &stdout, &stderr))
	generated, err := os.ReadFile(output)
	require.NoError(t, err)
	assert.Contains(t, string(generated), "//line ../dsl/flow.chain:")
	assert.NotContains(t, string(generated), dir)
}

func TestRunRejectsAliasedInputOutput(t *testing.T) {
	for _, tc := range []struct {
		name    string
		link    func(string, string) error
		reverse bool
	}{
		{"input symlink", os.Symlink, true},
		{"output symlink", os.Symlink, false},
		{"hard link", os.Link, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			original, alias := filepath.Join(dir, "flow.chain"), filepath.Join(dir, "alias.chain")
			const source = "package p\nworkflow Build(a Action) Flow[int] { nodes: a; success: a --> end }\n"
			require.NoError(t, os.WriteFile(original, []byte(source), 0644))
			require.NoError(t, tc.link(original, alias))
			input, output := original, alias
			if tc.reverse {
				input, output = alias, original
			}
			var stdout, stderr bytes.Buffer
			require.ErrorContains(t, run([]string{"-o", output, input}, &stdout, &stderr), "must differ")
			assert.Empty(t, stdout.Bytes())
			for _, path := range []string{original, alias} {
				contents, err := os.ReadFile(path)
				require.NoError(t, err)
				assert.Equal(t, source, string(contents), "source changed at %s", path)
			}
			originalInfo, err := os.Stat(original)
			require.NoError(t, err)
			aliasInfo, err := os.Stat(alias)
			require.NoError(t, err)
			assert.True(t, os.SameFile(originalInfo, aliasInfo), "input/output alias was replaced")
		})
	}
}
