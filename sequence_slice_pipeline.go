package chain

import (
	"context"
	"errors"
	"fmt"
	"github.com/JSYoo5B/chain/internal/logger"
	"runtime/debug"
)

func NewSequenceSlicePipeline[T any](name string, action Action[T], stopOnError bool) Action[[]T] {
	return &sequenceSlicePipeline[T]{
		name:        name,
		action:      action,
		stopOnError: stopOnError,
	}
}

type sequenceSlicePipeline[T any] struct {
	name        string
	action      Action[T]
	stopOnError bool
}

func (s sequenceSlicePipeline[T]) Name() string { return s.name }
func (s sequenceSlicePipeline[T]) Run(ctx context.Context, input []T) (output []T, err error) {
	pCtx := logger.WithRunnerDepth(ctx, s.name)
	output = make([]T, len(input))
	copy(output, input)

	// Wrap panic handling for safe running in a pipeline
	defer func() {
		if panicErr := recover(); panicErr != nil {
			logger.Errorf(pCtx, "chain: panic occurred on running, caused by %v", panicErr)
			debug.PrintStack()

			switch x := panicErr.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("unknown panic type")
			}
		}
	}()

	for i, in := range input {
		logger.Debugf(pCtx, "chain: running index %d", i)

		c := logger.WithRunnerDepth(ctx, fmt.Sprintf("%s[%d]/%s", s.name, i, s.action.Name()))

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
