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

// NewSimpleBranchAction creates a new BranchAction with customizable directions.
// It takes a name for the action, a slice of directions to define the possible control flow,
// and a function (runFunc) that defines the execution logic for the action.
func NewSimpleBranchAction[T any](name string, runFunc RunFunc[T], directions []string, branchFunc BranchFunc[T]) Action[T] {
	return &simpleBranchAction[T]{
		name:       name,
		runFunc:    runFunc,
		directions: directions,
		branchFunc: branchFunc,
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

type simpleBranchAction[T any] struct {
	name       string
	directions []string
	runFunc    RunFunc[T]
	branchFunc BranchFunc[T]
}

func (s simpleBranchAction[T]) Name() string         { return s.name }
func (s simpleBranchAction[T]) Directions() []string { return s.directions }
func (s simpleBranchAction[T]) Run(ctx context.Context, input T) (output T, err error) {
	if s.runFunc == nil {
		return input, nil
	}
	return s.runFunc(ctx, input)
}
func (s simpleBranchAction[T]) NextDirection(ctx context.Context, output T) (string, error) {
	if s.branchFunc == nil {
		return Success, nil
	}
	return s.branchFunc(ctx, output)
}
