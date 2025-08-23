package chain

import (
	"context"
	"fmt"
	internalErrors "github.com/JSYoo5B/chain/internal/errors"
	"github.com/JSYoo5B/chain/internal/logger"
	"runtime/debug"
	"sync"
)

// NewParallelSliceAction creates an Action that processes a slice's elements in parallel.
// Each element is transformed by the given action concurrently, maintaining the original order.
//
// The pipeline handles panics gracefully, continuing execution of other goroutines
// when one fails. If any error or panic occurs, the pipeline returns an error
// but still provides the processed output for successful operations.
func NewParallelSliceAction[T any](name string, action Action[T]) Action[[]T] {
	return &parallelSliceAction[T]{
		name:   name,
		action: action,
	}
}

type parallelSliceAction[T any] struct {
	name   string
	action Action[T]
}

func (p parallelSliceAction[T]) Name() string { return p.name }
func (p parallelSliceAction[T]) Run(ctx context.Context, input []T) (output []T, err error) {
	pCtx := logger.WithRunnerDepth(ctx, p.name)
	output = make([]T, len(input))
	copy(output, input)

	wg := sync.WaitGroup{}
	wg.Add(len(input))
	runIndex := func(i int, in T) {
		logger.Debugf(pCtx, "chain: running index %d", i)

		c := logger.WithRunnerDepth(ctx, fmt.Sprintf("%s[%d]/%s", p.name, i, p.action.Name()))
		runnerName, _ := logger.RunnerNameFromContext(c)

		// Wrap panic handling for safe running in a pipeline
		defer func() {
			if panicErr := recover(); panicErr != nil {
				logger.Errorf(pCtx, "chain: panic occurred on running index %d, caused by %v", i, panicErr)
				debug.PrintStack()

				output[i] = in
				err = internalErrors.NewPanicError(runnerName, panicErr)
				wg.Done()
				return
			}
		}()

		out, e := p.action.Run(c, in)
		if e != nil {
			logger.Errorf(pCtx, "chain: error occurred in index %d: %v", i, e)
			err = e
		}
		output[i] = out
		wg.Done()
	}
	for i, in := range input {
		go runIndex(i, in)
	}
	wg.Wait()

	return output, err
}
