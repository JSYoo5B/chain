package chain

import (
	"context"
	"fmt"
	internalErrors "github.com/JSYoo5B/chain/internal/errors"
	"github.com/JSYoo5B/chain/internal/logger"
	"runtime/debug"
)

// NewSequenceSliceAction creates an Action that processes a slice's elements sequentially.
// Each element is transformed by the given action one at a time, maintaining the original order.
//
// The stopOnError parameter controls error handling behavior:
// - When true: stops processing immediately on the first error, leaving remaining elements unchanged
// - When false: continues processing all elements even if errors occur
// Panics always stop execution regardless of the stopOnError setting.
func NewSequenceSliceAction[T any](name string, action Action[T], stopOnError bool) Action[[]T] {
	if action == nil {
		panic("action cannot be nil")
	}

	return &sequenceSliceAction[T]{
		name:        name,
		action:      action,
		stopOnError: stopOnError,
	}
}

type sequenceSliceAction[T any] struct {
	name        string
	action      Action[T]
	stopOnError bool
}

func (s sequenceSliceAction[T]) Name() string { return s.name }
func (s sequenceSliceAction[T]) Run(ctx context.Context, input []T) (output []T, err error) {
	pCtx := logger.WithRunnerDepth(ctx, s.name)
	runnerName, _ := logger.RunnerNameFromContext(pCtx)
	output = make([]T, len(input))
	copy(output, input)

	// Wrap panic handling for safe running in an action
	defer func() {
		if panicErr := recover(); panicErr != nil {
			logger.Errorf(pCtx, "chain: panic occurred on running, caused by %v", panicErr)
			debug.PrintStack()

			err = internalErrors.NewPanicError(runnerName, panicErr)
		}
	}()

	for i, in := range input {
		logger.Debugf(pCtx, "chain: running index %d", i)

		c := logger.WithRunnerDepth(ctx, fmt.Sprintf("%s[%d]/%s", s.name, i, s.action.Name()))
		runnerName, _ = logger.RunnerNameFromContext(c)

		out, e := s.action.Run(c, in)
		output[i] = out
		if e != nil {
			if s.stopOnError {
				logger.Errorf(pCtx, "chain: stopping after error in index %d", i)
				return output, e

			} else {
				logger.Errorf(pCtx, "chain: error occurred in index %d: %v", i, e)
				err = e
			}
		}
	}

	return output, err
}
