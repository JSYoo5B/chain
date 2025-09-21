package parser

import (
	"github.com/antlr4-go/antlr/v4"
)

type DiagnosticListener struct {
	*antlr.DefaultErrorListener
	SyntaxErrors        []SyntaxError
	AmbiguityWarnings   []AmbiguityWarning
	PerformanceWarnings []PerformanceWarning
}

func NewDiagnosticListener() *DiagnosticListener {
	return &DiagnosticListener{
		DefaultErrorListener: antlr.NewDefaultErrorListener(),
		SyntaxErrors:         []SyntaxError{},
		AmbiguityWarnings:    []AmbiguityWarning{},
		PerformanceWarnings:  []PerformanceWarning{},
	}
}

func (l *DiagnosticListener) SyntaxError(_ antlr.Recognizer, _ interface{}, line, column int, msg string, _ antlr.RecognitionException) {
	err := SyntaxError{
		Line:    line,
		Column:  column,
		Message: msg,
	}
	l.SyntaxErrors = append(l.SyntaxErrors, err)
}

func (l *DiagnosticListener) HasSyntaxError() bool {
	return len(l.SyntaxErrors) > 0
}

func (l *DiagnosticListener) ReportAmbiguity(recognizer antlr.Parser, _ *antlr.DFA, startIndex, stopIndex int, _ bool, ambigAlts *antlr.BitSet, configs *antlr.ATNConfigSet) {
	tokens := recognizer.GetTokenStream()
	input := tokens.GetTextFromInterval(antlr.NewInterval(startIndex, stopIndex))
	parserRuleCtx := recognizer.GetParserRuleContext()

	warning := AmbiguityWarning{
		StartIndex:   startIndex,
		StopIndex:    stopIndex,
		Input:        input,
		Alternatives: ambigAlts.String(),
		RuleStack:    recognizer.GetRuleInvocationStack(parserRuleCtx),
		Configs:      configs.String(),
	}
	l.AmbiguityWarnings = append(l.AmbiguityWarnings, warning)
}

func (l *DiagnosticListener) HasAmbiguityWarnings() bool {
	return len(l.AmbiguityWarnings) > 0
}

func (l *DiagnosticListener) ReportAttemptingFullContext(recognizer antlr.Parser, _ *antlr.DFA, startIndex, stopIndex int, _ *antlr.BitSet, _ *antlr.ATNConfigSet) {
	tokens := recognizer.GetTokenStream()
	input := tokens.GetTextFromInterval(antlr.NewInterval(startIndex, stopIndex))
	parserRuleCtx := recognizer.GetParserRuleContext()

	warning := PerformanceWarning{
		StartIndex: startIndex,
		StopIndex:  stopIndex,
		Input:      input,
		RuleStack:  recognizer.GetRuleInvocationStack(parserRuleCtx),
	}
	l.PerformanceWarnings = append(l.PerformanceWarnings, warning)
}

func (l *DiagnosticListener) ReportContextSensitivity(recognizer antlr.Parser, _ *antlr.DFA, startIndex, stopIndex, _ int, _ *antlr.ATNConfigSet) {
	tokens := recognizer.GetTokenStream()
	input := tokens.GetTextFromInterval(antlr.NewInterval(startIndex, stopIndex))
	parserRuleCtx := recognizer.GetParserRuleContext()

	warning := PerformanceWarning{
		StartIndex: startIndex,
		StopIndex:  stopIndex,
		Input:      input,
		RuleStack:  recognizer.GetRuleInvocationStack(parserRuleCtx),
	}
	l.PerformanceWarnings = append(l.PerformanceWarnings, warning)
}

func (l *DiagnosticListener) HasPerformanceWarnings() bool {
	return len(l.PerformanceWarnings) > 0
}
