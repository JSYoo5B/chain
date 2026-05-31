package parser

import "fmt"

type SyntaxError struct {
	Line    int
	Column  int
	Message string
}

func (e *SyntaxError) String() string {
	return fmt.Sprintf("%d:%d - %s", e.Line, e.Column, e.Message)
}

type AmbiguityWarning struct {
	StartIndex   int
	StopIndex    int
	Input        string
	Alternatives string
	RuleStack    []string
	Configs      string
}

func (w *AmbiguityWarning) String() string {
	return fmt.Sprintf(
		"In rule stack %v, input `%s` could be interpreted as alternatives: %s",
		w.RuleStack, w.Input, w.Alternatives)
}

type PerformanceWarning struct {
	StartIndex int
	StopIndex  int
	Input      string
	RuleStack  []string
}

func (w *PerformanceWarning) String() string {
	return fmt.Sprintf(
		"In rule stack %v, parser struggled with input `%s`. Grammar may be inefficient",
		w.RuleStack, w.Input)
}
