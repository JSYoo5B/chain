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
