package chain

import (
	"context"
	"fmt"
	internalErrors "github.com/JSYoo5B/chain/internal/errors"
	"github.com/JSYoo5B/chain/internal/logger"
	"maps"
	"runtime/debug"
)

// NewSequenceMapPipeline creates an Action that processes a map's values sequentially.
// Each value is transformed by the given action one at a time, maintaining the original keys.
//
// Unlike parallel processing, sequential execution stops immediately when a panic occurs,
// leaving unprocessed values unchanged in the output.
func NewSequenceMapPipeline[K comparable, T any](name string, action Action[T]) Action[map[K]T] {
	return &sequenceMapPipeline[K, T]{
		name:   name,
		action: action,
	}
}

type sequenceMapPipeline[K comparable, T any] struct {
	name   string
	action Action[T]
}

func (s sequenceMapPipeline[K, T]) Name() string { return s.name }
func (s sequenceMapPipeline[K, T]) Run(ctx context.Context, input map[K]T) (output map[K]T, err error) {
	pCtx := logger.WithRunnerDepth(ctx, s.name)
	runnerName, _ := logger.RunnerNameFromContext(pCtx)
	output = make(map[K]T)
	maps.Copy(output, input)

	// Wrap panic handling for safe running in a pipeline
	defer func() {
		if panicErr := recover(); panicErr != nil {
			logger.Errorf(pCtx, "chain: panic occurred on running, caused by %v", panicErr)
			debug.PrintStack()

			err = internalErrors.NewPanicError(runnerName, panicErr)
		}
	}()

	for k, in := range input {
		logger.Debugf(pCtx, "chain: running key `%v`", k)

		c := logger.WithRunnerDepth(ctx, fmt.Sprintf("%s[%v]/%s", s.name, k, s.action.Name()))
		runnerName, _ = logger.RunnerNameFromContext(c)

		out, e := s.action.Run(c, in)
		if e != nil {
			logger.Errorf(pCtx, "chain: error occurred in key `%v`: %v", k, e)
			err = e
		}
		output[k] = out
	}

	return output, err
}
