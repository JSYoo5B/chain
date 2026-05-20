package chain

import "context"

// RunFunc is the function signature used by NewSimpleAction.
// It receives an input value and returns the next value or an error.
type RunFunc[T any] func(ctx context.Context, input T) (output T, err error)

// NewSimpleAction creates a new Action with a custom Run function,
// which can be a pure function or closure.
// The provided runFunc must match the RunFunc signature.
//
// This allows simple Actions to be created without manually defining a struct.
func NewSimpleAction[T any](name string, runFunc RunFunc[T]) Action[T] {
	if runFunc == nil {
		panic("runFunc cannot be nil")
	}

	return &simpleAction[T]{
		name:    name,
		runFunc: runFunc,
	}
}

type simpleAction[T any] struct {
	name    string
	runFunc RunFunc[T]
}

func (s simpleAction[T]) Name() string { return s.name }
func (s simpleAction[T]) Run(ctx context.Context, input T) (output T, err error) {
	return s.runFunc(ctx, input)
}
