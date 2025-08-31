package chain

import "context"

// RunFunc defines the signature of a function used to implement an Action's execution logic.
// It is a function that takes an input of type T (a generic type) and returns an output of type T
// along with any error encountered during execution.
type RunFunc[T any] func(ctx context.Context, input T) (output T, err error)

// NewSimpleAction creates a new Action with a custom Run function,
// which can be a pure function or closure.
// The provided runFunc must match the RunFunc signature, where T is a generic type representing
// the input and output types for the Action's execution.
//
// This allows for the creation of simple Actions without manually defining a separate struct
// that implements the Action interface.
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
