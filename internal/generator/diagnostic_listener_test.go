package generator

import (
	"embed"
	"github.com/JSYoo5B/chain/internal/parser"
	"github.com/antlr4-go/antlr/v4"
	"io/fs"
	"testing"
)

//go:embed testdata/*.chain
var testDslFs embed.FS

func TestChainParserAmbiguity(t *testing.T) {
	testCases := make(map[string]string)

	fs.WalkDir(testDslFs, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		data, err := testDslFs.ReadFile(path)
		if err != nil {
			return err
		}
		testCases[path] = string(data)
		return nil
	})

	for name, source := range testCases {
		t.Run(name, func(t *testing.T) {
			inputStream := antlr.NewInputStream(source)
			lexer := parser.NewCommonLexer(inputStream)
			stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
			p := parser.NewChainParser(stream)

			p.RemoveErrorListeners()

			listener := newAmbiguityListener()
			p.AddErrorListener(listener)

			p.GetInterpreter().SetPredictionMode(antlr.PredictionModeLLExactAmbigDetection)

			p.SourceFile()

			if listener.hasAmbiguity() {
				t.Logf("Grammar ambiguity detected in source:\n---INPUT---\n%s\n---AMBIGUITIES---\n%s\n-----------", source, listener.getAmbiguity())
			}
			if listener.hasSyntaxError() {
				t.Errorf("Grammar syntax error detected in source:\n---INPUT---\n%s\n---SYNTAX ERRORS---\n%s\n-----------", source, listener.getSyntaxError())
			}
		})
	}
}
