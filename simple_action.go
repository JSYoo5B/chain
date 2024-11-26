package railway

import "context"

// NewSimpleAction creates a new Action with a fixed set of default directions.
// It takes a name and a function (runFunc) that defines the execution logic for the action.
func NewSimpleAction[T any](name string, runFunc RunFunc[T]) Action[T] {
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
	if s.runFunc == nil {
		return input, nil
	}
	return s.runFunc(ctx, input)
}
