package golang

import (
	"github.com/JSYoo5B/chain/internal/dsl/parser"
	"github.com/antlr4-go/antlr/v4"
)

type SimpleDiagnosticListener struct {
	*antlr.DefaultErrorListener
	SyntaxErrors        []parser.SyntaxError
	AmbiguityWarnings   []parser.AmbiguityWarning
	PerformanceWarnings []parser.PerformanceWarning
}

func NewSimpleDiagnosticListener() *SimpleDiagnosticListener {
	return &SimpleDiagnosticListener{
		DefaultErrorListener: antlr.NewDefaultErrorListener(),
		SyntaxErrors:         []parser.SyntaxError{},
		AmbiguityWarnings:    []parser.AmbiguityWarning{},
		PerformanceWarnings:  []parser.PerformanceWarning{},
	}
}

func (l *SimpleDiagnosticListener) SyntaxError(_ antlr.Recognizer, _ interface{}, line, column int, msg string, _ antlr.RecognitionException) {
	err := parser.SyntaxError{
		Line:    line,
		Column:  column,
		Message: msg,
	}
	l.SyntaxErrors = append(l.SyntaxErrors, err)
}

func (l *SimpleDiagnosticListener) HasSyntaxError() bool {
	return len(l.SyntaxErrors) > 0
}

func (l *SimpleDiagnosticListener) ReportAmbiguity(recognizer antlr.Parser, _ *antlr.DFA, startIndex, stopIndex int, _ bool, ambigAlts *antlr.BitSet, configs *antlr.ATNConfigSet) {
	tokens := recognizer.GetTokenStream()
	input := tokens.GetTextFromInterval(antlr.NewInterval(startIndex, stopIndex))
	parserRuleCtx := recognizer.GetParserRuleContext()

	warning := parser.AmbiguityWarning{
		StartIndex:   startIndex,
		StopIndex:    stopIndex,
		Input:        input,
		Alternatives: ambigAlts.String(),
		RuleStack:    recognizer.GetRuleInvocationStack(parserRuleCtx),
		Configs:      configs.String(),
	}
	l.AmbiguityWarnings = append(l.AmbiguityWarnings, warning)
}

func (l *SimpleDiagnosticListener) HasAmbiguityWarnings() bool {
	return len(l.AmbiguityWarnings) > 0
}

func (l *SimpleDiagnosticListener) ReportAttemptingFullContext(recognizer antlr.Parser, _ *antlr.DFA, startIndex, stopIndex int, _ *antlr.BitSet, _ *antlr.ATNConfigSet) {
	tokens := recognizer.GetTokenStream()
	input := tokens.GetTextFromInterval(antlr.NewInterval(startIndex, stopIndex))
	parserRuleCtx := recognizer.GetParserRuleContext()

	warning := parser.PerformanceWarning{
		StartIndex: startIndex,
		StopIndex:  stopIndex,
		Input:      input,
		RuleStack:  recognizer.GetRuleInvocationStack(parserRuleCtx),
	}
	l.PerformanceWarnings = append(l.PerformanceWarnings, warning)
}

func (l *SimpleDiagnosticListener) ReportContextSensitivity(recognizer antlr.Parser, _ *antlr.DFA, startIndex, stopIndex, _ int, _ *antlr.ATNConfigSet) {
	tokens := recognizer.GetTokenStream()
	input := tokens.GetTextFromInterval(antlr.NewInterval(startIndex, stopIndex))
	parserRuleCtx := recognizer.GetParserRuleContext()

	warning := parser.PerformanceWarning{
		StartIndex: startIndex,
		StopIndex:  stopIndex,
		Input:      input,
		RuleStack:  recognizer.GetRuleInvocationStack(parserRuleCtx),
	}
	l.PerformanceWarnings = append(l.PerformanceWarnings, warning)
}

func (l *SimpleDiagnosticListener) HasPerformanceWarnings() bool {
	return len(l.PerformanceWarnings) > 0
}
