package chain

import "context"

// BranchFunc represents the signature for the function that defines the branching logic
// for a BranchAction in the package. It takes the running context and output as input
// and returns the direction for the next step in the process along with any potential error.
type BranchFunc[T any] func(ctx context.Context, output T) (direction string, err error)

// NewSimpleBranchAction creates a new BranchAction with customizable directions.
// It accepts a name for the action, a slice of directions that define the possible control flow,
// and a BranchFunc that contains the branching logic, which dictates the next direction
// based on the action's output.
//
// Additionally, a custom runFunc can be provided to define the execution logic of the action.
// This function must match the RunFunc signature, where T is the generic type representing
// the input and output types for the action. If no specific execution logic is needed, the
// runFunc can be provided as `nil`. In this case, the action will simply pass the input
// through to the output without modification.
//
// This allows for the creation of simple BranchActions without manually defining a separate struct
// that implements the BranchAction interface.
func NewSimpleBranchAction[T any](name string, runFunc RunFunc[T], directions []string, branchFunc BranchFunc[T]) BranchAction[T] {
	if len(directions) == 0 {
		panic("directions cannot be empty")
	} else if branchFunc == nil {
		panic("branchFunc cannot be nil")
	}

	if runFunc == nil {
		runFunc = func(_ context.Context, input T) (T, error) { return input, nil }
	}
	return &simpleBranchAction[T]{
		name:       name,
		runFunc:    runFunc,
		directions: directions,
		branchFunc: branchFunc,
	}
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
	return s.runFunc(ctx, input)
}
func (s simpleBranchAction[T]) NextDirection(ctx context.Context, output T) (string, error) {
	return s.branchFunc(ctx, output)
}
