package golang

import (
	"embed"
	"github.com/antlr4-go/antlr/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/fs"
	"strings"
	"testing"
)

//go:embed testdata/antlrExamples/*.go_
var testAntlrExampleFs embed.FS

func TestAntlrExamplesHasError(t *testing.T) {
	testCases := make(map[string]string)

	err := fs.WalkDir(testAntlrExampleFs, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		} else if d.IsDir() {
			return nil
		}

		data, err := testAntlrExampleFs.ReadFile(path)
		if err != nil {
			return err
		}
		testName := strings.TrimPrefix(path, "testdata/antlrExamples/")
		testName = strings.TrimSuffix(testName, ".go_")
		testCases[testName] = string(data)
		return nil
	})
	require.NoError(t, err)

	for name, source := range testCases {
		t.Run(name, func(t *testing.T) {
			inputStream := antlr.NewInputStream(source)
			lexer := NewCommonLexer(inputStream)
			stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
			p := NewGoParser(stream)
			p.RemoveErrorListeners()
			listener := NewSimpleDiagnosticListener()
			p.AddErrorListener(listener)
			p.GetInterpreter().SetPredictionMode(antlr.PredictionModeLLExactAmbigDetection)
			p.SourceFile()

			t.Run("syntax error", func(t *testing.T) {
				assert.False(t, listener.HasSyntaxError())
				for _, syntaxErr := range listener.SyntaxErrors {
					t.Log(syntaxErr)
				}
			})
			t.Run("ambiguity", func(t *testing.T) {
				if !listener.HasAmbiguityWarnings() {
					return
				}
				for _, ambigWarn := range listener.AmbiguityWarnings {
					t.Log(ambigWarn.String())
				}
			})
			t.Run("performance", func(t *testing.T) {
				if !listener.HasPerformanceWarnings() {
					return
				}
				for _, perfWarn := range listener.PerformanceWarnings {
					t.Log(perfWarn.String())
				}
			})
		})
	}
}
