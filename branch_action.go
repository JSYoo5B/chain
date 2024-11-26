package chain

import "context"

// BranchAction is an interface for actions that control branching in the execution flow
// of a Pipeline. It extends the Action interface and adds methods for handling conditional
// branching based on the execution results.
type BranchAction[T any] interface {
	// Name returns the name of the BranchAction.
	Name() string

	// Run executes the branch action, optionally modifying the input and returning an output.
	// If the input doesn't need changes, it can be passed through as output. The method also
	// returns an error if the action cannot be executed successfully.
	Run(ctx context.Context, input T) (output T, err error)

	// Directions returns a list of possible directions that the pipeline can take.
	// These directions are used for validation and must include all possible values that
	// NextDirection can return.
	Directions() []string

	// NextDirection determines the next execution path based on the result of Run.
	// It is called only if Run succeeds (err == nil).
	// The method returns a direction from the list defined by Directions.
	NextDirection(ctx context.Context, output T) (direction string, err error)
}
