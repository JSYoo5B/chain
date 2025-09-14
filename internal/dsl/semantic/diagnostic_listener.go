package semantic

import (
	"fmt"
	"github.com/antlr4-go/antlr/v4"
	"strings"
	"sync"
)

type diagnosticListener struct {
	*antlr.DefaultErrorListener
	ambiguities  []string
	syntaxErrors []string
	mutex        sync.Mutex
}

func newAmbiguityListener() *diagnosticListener {
	return new(diagnosticListener)
}

func (l *diagnosticListener) ReportAmbiguity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, exact bool, ambigAlts *antlr.BitSet, configs *antlr.ATNConfigSet) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	tokens := recognizer.GetTokenStream()
	input := tokens.GetTextFromInterval(antlr.NewInterval(startIndex, stopIndex))

	msg := fmt.Sprintf("Ambiguity found at input '%s', alternatives: %s", input, ambigAlts.String())
	l.ambiguities = append(l.ambiguities, msg)
}

func (l *diagnosticListener) hasAmbiguity() bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	return len(l.ambiguities) > 0
}

func (l *diagnosticListener) getAmbiguity() string {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	return strings.Join(l.ambiguities, "\n")
}

func (l *diagnosticListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	errorMsg := fmt.Sprintf("Syntax error at line %d:%d - %s", line, column, msg)
	l.syntaxErrors = append(l.syntaxErrors, errorMsg)
}

func (l *diagnosticListener) hasSyntaxError() bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	return len(l.syntaxErrors) > 0
}

func (l *diagnosticListener) getSyntaxError() string {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	return strings.Join(l.syntaxErrors, "\n")
}
