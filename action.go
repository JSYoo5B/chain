package railway

import (
	"context"
)

// RunFunc represents the signature for the function that defines the execution logic in the railway.
// For more details, refer to the Action.Run method.
type RunFunc[T any] func(ctx context.Context, input T) (output T, err error)

type Action[T any] interface {
	// Name returns the name of the Action.
	// This is typically a unique identifier for the Action that can be
	// used to distinguish it from other actions in the Pipeline.
	Name() string

	// Run executes the Action with the given context and input, and returns two values:
	// - output: The result of the Action's execution.
	// - err: An error indicating whether the Action failed during execution.
	Run(ctx context.Context, input T) (output T, err error)
}

// Terminate returns a termination action,
// providing a clear indication of termination rather than using nil.
func Terminate[T any]() Action[T] {
	return nil
}
