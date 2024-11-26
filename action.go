package railway

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

// Terminate explicitly ends execution in a Pipeline by returning nil.
// It signals that no further actions will be executed.
//
// Use in ActionPlan to clearly indicate termination intent.
func Terminate[T any]() Action[T] {
	return nil
}
