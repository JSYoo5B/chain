package railway

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"runtime/debug"
)

// RunFunc represents the signature for the function that defines the execution logic in the railway.
// For more details, refer to the Action.Run method.
type RunFunc[T any] func(ctx context.Context, input T) (output T, direction string, err error)

type Action[T any] interface {
	// Name returns the name of the action.
	// This is typically a unique identifier for the action that can be
	// used to distinguish it from other actions in the pipeline.
	Name() string
	// Directions returns the possible directions for the action to take as a slice of strings.
	// These directions (e.g., Success, Error, Abort) define how the action can
	// proceed based on the outcome of its execution and guide the flow in the Pipeline.
	// Additional custom directions can also be provided, which can be used to implement custom
	// branching logic and control the flow of execution in the Pipeline beyond the default flow.
	// Default directions (Success, Error, Abort) can be omitted.
	Directions() []string
	// Run executes the action with the given context and input, and returns three values:
	// - output: The result of the Action's execution.
	// - direction: A string indicating the flow direction after the action completes.
	// - err: An error indicating whether the action failed during execution.
	Run(ctx context.Context, input T) (output T, direction string, err error)
}

// Terminate returns a termination action,
// providing a clear indication of termination rather than returning nil.
func Terminate[T any]() Action[T] {
	return nil
}

func runAction[T any](action Action[T], ctx context.Context, input T) (output T, direction string, runError error) {
	// Wrap panic handling for safe running in pipeline
	defer func() {
		if panicErr := recover(); panicErr != nil {
			logrus.Errorf("%s: panic occurred on running, caused by %s", action.Name(), panicErr)
			debug.PrintStack()

			output = input
			direction = Abort
			switch x := panicErr.(type) {
			case string:
				runError = errors.New(x)
			case error:
				runError = x
			default:
				runError = errors.New("unknown panic type")
			}
		}
	}()

	output, direction, runError = action.Run(ctx, input)
	if runError != nil && direction != Error && direction != Abort {
		logrus.Errorf("%s: invoked error but directing `%s`, overriding with `error`", action.Name(), direction)
		direction = Error
	}
	return output, direction, runError
}
