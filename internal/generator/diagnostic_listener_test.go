package generator

import (
	_ "embed"
	"github.com/JSYoo5B/chain/internal/parser"
	"github.com/antlr4-go/antlr/v4"
	"testing"
)

//go:embed testdata/helloworld.chain
var helloWorld string

//go:embed testdata/collatz.chain
var collatz string

func TestChainParserAmbiguity(t *testing.T) {
	testCases := map[string]struct {
		input string
	}{
		"hello world": {
			input: helloWorld,
		},
		"collatz": {
			input: collatz,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			inputStream := antlr.NewInputStream(tc.input)
			lexer := parser.NewCommonLexer(inputStream)
			stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
			p := parser.NewChainParser(stream)

			p.RemoveErrorListeners()

			listener := newAmbiguityListener()
			p.AddErrorListener(listener)

			p.GetInterpreter().SetPredictionMode(antlr.PredictionModeLLExactAmbigDetection)

			p.SourceFile()

			if listener.hasAmbiguity() {
				t.Errorf("Grammar ambiguity detected in input:\n---INPUT---\n%s\n---AMBIGUITIES---\n%s\n-----------", tc.input, listener.getAmbiguity())
			}
			if listener.hasSyntaxError() {
				t.Errorf("Grammar syntax error detected in input:\n---INPUT---\n%s\n---SYNTAX ERRORS---\n%s\n-----------", tc.input, listener.getSyntaxError())
			}
		})
	}
}
