package chain

import (
	"context"
)

// Action is the basic unit of execution in a package.
// It represents a single task that processes input and produces output.
type Action[T any] interface {
	// Name provides the identifier of this Action.
	Name() string

	// Run executes the Action, processing the input and returning output or an error.
	Run(ctx context.Context, input T) (output T, err error)
}

// BranchAction is an interface for actions that control branching in the execution flow
// of a Workflow. It extends the Action interface and adds methods for handling conditional
// branching based on the execution results.
type BranchAction[T any] interface {
	// Name returns the name of the BranchAction.
	Name() string

	// Run executes the branch action, optionally modifying the input and returning an output.
	// If the input doesn't need changes, it can be passed through as output. The method also
	// returns an error if the action cannot be executed successfully.
	Run(ctx context.Context, input T) (output T, err error)

	// Directions return a list of possible directions that the Workflow can take.
	// These directions are used for validation and must include all possible values that
	// NextDirection can return.
	Directions() []string

	// NextDirection determines the next execution path based on the result of Run.
	// It is called only if Run succeeds (err == nil).
	// The method returns a direction from the list defined by Directions.
	NextDirection(ctx context.Context, output T) (direction string, err error)
}

// Terminate explicitly ends execution in a Workflow by returning nil.
// It signals that no further actions will be executed.
//
// Use in ActionPlan to clearly indicate termination intent.
func Terminate[T any]() Action[T] {
	return nil
}
