package railway

import "context"

// BranchFunc represents the signature for the function that defines the branching logic in the railway.
// For more details, refer to the BranchAction.NextDirection method.
type BranchFunc[T any] func(ctx context.Context, output T) (direction string, err error)

type BranchAction[T any] interface {
	// Name returns the name of the BranchAction.
	// This is typically a unique identifier for the BranchAction that can be
	// used to distinguish it from other actions in the Pipeline.
	Name() string

	// Run executes the BranchAction with the given context and input, and returns two values:
	// - output: The result of the Action's execution.
	// - err: An error indicating whether the Action failed during execution.
	Run(ctx context.Context, input T) (output T, err error)

	// Directions returns the possible directions for branching as a slice of strings.
	// These directions define how the BranchAction can proceed based on the outcome of
	// its execution and guide the flow in the Pipeline.
	Directions() []string

	NextDirection(ctx context.Context, output T) (string, error)
}
