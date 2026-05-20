package chain

import "context"

// BranchFunc represents the signature for the function that defines the branching logic
// for a BranchAction in the package. It takes the running context and output as input
// and returns the direction for the next step along with any error.
type BranchFunc[T any] func(ctx context.Context, output T) (direction string, err error)

// NewSimpleBranchAction creates a new BranchAction with customizable directions.
// It accepts a name, custom directions, and a BranchFunc that selects the next
// direction from the action's output.
//
// A custom runFunc can be provided to define execution logic. If runFunc is nil,
// the action passes the input through as output.
//
// This allows simple BranchActions to be created without manually defining a
// struct that implements BranchAction.
func NewSimpleBranchAction[T any](name string, runFunc RunFunc[T], directions []string, branchFunc BranchFunc[T]) BranchAction[T] {
	if len(directions) == 0 {
		panic("directions cannot be empty")
	} else if branchFunc == nil {
		panic("branchFunc cannot be nil")
	}

	if runFunc == nil {
		runFunc = func(_ context.Context, input T) (T, error) { return input, nil }
	}
	copiedDirections := make([]string, len(directions))
	copy(copiedDirections, directions)

	return &simpleBranchAction[T]{
		name:       name,
		runFunc:    runFunc,
		directions: copiedDirections,
		branchFunc: branchFunc,
	}
}

type simpleBranchAction[T any] struct {
	name       string
	directions []string
	runFunc    RunFunc[T]
	branchFunc BranchFunc[T]
}

func (s simpleBranchAction[T]) Name() string { return s.name }
func (s simpleBranchAction[T]) Directions() []string {
	directions := make([]string, len(s.directions))
	copy(directions, s.directions)
	return directions
}
func (s simpleBranchAction[T]) Run(ctx context.Context, input T) (output T, err error) {
	return s.runFunc(ctx, input)
}
func (s simpleBranchAction[T]) NextDirection(ctx context.Context, output T) (string, error) {
	return s.branchFunc(ctx, output)
}
