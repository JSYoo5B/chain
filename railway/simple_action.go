package railway

import "context"

// NewSimpleAction creates a new Action with a fixed set of default directions.
// It takes a name and a function (runFunc) that defines the execution logic for the action.
func NewSimpleAction[T any](name string, runFunc RunFunc[T]) Action[T] {
	return &simpleAction[T]{
		name:       name,
		directions: []string{Success, Error, Abort},
		runFunc:    runFunc,
	}
}

// NewSimpleBranchAction creates a new Action with customizable directions.
// It takes a name for the action, a slice of directions to define the possible control flow,
// and a function (runFunc) that defines the execution logic for the action.
func NewSimpleBranchAction[T any](name string, directions []string, runFunc RunFunc[T]) Action[T] {
	return &simpleAction[T]{
		name:       name,
		directions: directions,
		runFunc:    runFunc,
	}
}

type simpleAction[T any] struct {
	name       string
	directions []string
	runFunc    RunFunc[T]
}

func (s simpleAction[T]) Name() string         { return s.name }
func (s simpleAction[T]) Directions() []string { return s.directions }
func (s simpleAction[T]) Run(ctx context.Context, input T) (output T, direction string, err error) {
	return s.runFunc(ctx, input)
}
