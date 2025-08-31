package chain

import (
	"context"
	"fmt"
	internalErrors "github.com/JSYoo5B/chain/internal/errors"
	"github.com/JSYoo5B/chain/internal/logger"
	"maps"
	"runtime/debug"
	"sync"
)

// NewParallelMapAction creates an Action that processes a map's values in parallel.
// Each value is transformed by the given action concurrently, maintaining the original keys.
//
// The action handles panics gracefully, continuing execution of other goroutines
// when one fails. If any error or panic occurs, the action returns an error
// but still provides the processed output for successful operations.
func NewParallelMapAction[K comparable, T any](name string, action Action[T]) Action[map[K]T] {
	if action == nil {
		panic("action cannot be nil")
	}

	return &parallelMapAction[K, T]{
		name:   name,
		action: action,
	}
}

type parallelMapAction[K comparable, T any] struct {
	name   string
	action Action[T]
}

func (p parallelMapAction[K, T]) Name() string { return p.name }
func (p parallelMapAction[K, T]) Run(ctx context.Context, input map[K]T) (output map[K]T, err error) {
	pCtx := logger.WithRunnerDepth(ctx, p.name)
	output = make(map[K]T)
	maps.Copy(output, input)

	wg := sync.WaitGroup{}
	wg.Add(len(input))
	runKey := func(k K, in T) {
		logger.Debugf(pCtx, "chain: running key `%v`", k)
		c := logger.WithRunnerDepth(ctx, fmt.Sprintf("%s[%v]/%s", p.name, k, p.action.Name()))
		runnerName, _ := logger.RunnerNameFromContext(c)

		// Wrap panic handling for safe running in an action
		defer func() {
			if panicErr := recover(); panicErr != nil {
				logger.Errorf(pCtx, "chain: panic occurred on running key %v, caused by %v", k, panicErr)
				debug.PrintStack()

				output[k] = in
				err = internalErrors.NewPanicError(runnerName, panicErr)
				wg.Done()
				return
			}
		}()

		out, e := p.action.Run(c, in)
		if e != nil {
			logger.Errorf(pCtx, "chain: error occurred in key `%v`: %v", k, e)
			err = e
		}
		output[k] = out
		wg.Done()
	}
	for k, in := range input {
		go runKey(k, in)
	}
	wg.Wait()

	return output, err
}
