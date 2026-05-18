package chain

import (
	"context"
	"errors"
	"fmt"
	internalErrors "github.com/JSYoo5B/chain/internal/errors"
	"github.com/JSYoo5B/chain/internal/logger"
	"runtime/debug"
)

// AsParallelSliceAction creates an Action that processes a slice's elements in parallel.
// Each element is transformed by the given action concurrently, maintaining the original order.
//
// The action handles panics gracefully, continuing execution of other goroutines
// when one fails. If any error or panic occurs, the action returns an error
// but still provides the processed output for successful operations.
func AsParallelSliceAction[T any](name string, action Action[T]) Action[[]T] {
	if action == nil {
		panic("action cannot be nil")
	}

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

	type result struct {
		index  int
		output T
		err    error
	}

	results := make(chan result, len(input))
	runIndex := func(i int, in T) {
		logger.Debugf(pCtx, "chain: running index %d", i)

		c := logger.WithRunnerDepth(ctx, fmt.Sprintf("%s[%d]/%s", p.name, i, p.action.Name()))
		runnerName, _ := logger.RunnerNameFromContext(c)

		// Wrap panic handling for safe running in an action
		defer func() {
			if panicErr := recover(); panicErr != nil {
				logger.Errorf(pCtx, "chain: panic occurred on running index %d, caused by %v", i, panicErr)
				debug.PrintStack()

				results <- result{
					index:  i,
					output: in,
					err:    internalErrors.NewPanicError(runnerName, panicErr),
				}
				return
			}
		}()

		out, e := p.action.Run(c, in)
		if e != nil {
			logger.Errorf(pCtx, "chain: error occurred in index %d: %v", i, e)
		}
		results <- result{
			index:  i,
			output: out,
			err:    e,
		}
	}
	for i, in := range input {
		go runIndex(i, in)
	}

	for range len(input) {
		r := <-results
		output[r.index] = r.output
		err = errors.Join(err, r.err)
	}

	return output, err
}
