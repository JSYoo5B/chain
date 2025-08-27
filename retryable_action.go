package chain

import (
	"context"
	"fmt"
	internalErrors "github.com/JSYoo5B/chain/internal/errors"
	"github.com/JSYoo5B/chain/internal/logger"
	"runtime/debug"
)

func NewRetryableAction[T any](name string, mainAction, rollbackAction Action[T], maxRetry int) Action[T] {
	return &retryableAction[T]{
		name:           name,
		mainAction:     mainAction,
		rollbackAction: rollbackAction,
		maxRetry:       maxRetry,
	}
}

type retryableAction[T any] struct {
	name           string
	mainAction     Action[T]
	rollbackAction Action[T]
	maxRetry       int
}

func (r retryableAction[T]) Name() string { return r.name }
func (r retryableAction[T]) Run(ctx context.Context, input T) (output T, err error) {
	pCtx := logger.WithRunnerDepth(ctx, r.name)
	runnerName, _ := logger.RunnerNameFromContext(pCtx)
	output = input

	defer func() {
		if panicErr := recover(); panicErr != nil {
			logger.Errorf(pCtx, "chain: panic occurred on running, caused by %v", panicErr)
			debug.PrintStack()

			err = internalErrors.NewPanicError(runnerName, panicErr)
		}
	}()

	mCtx := logger.WithRunnerDepth(pCtx, r.mainAction.Name())
	rCtx := logger.WithRunnerDepth(pCtx, r.rollbackAction.Name())
	for attempt := 1; attempt <= r.maxRetry; attempt++ {
		logger.Debugf(pCtx, "chain: running %s, attempt: %d", r.mainAction.Name(), attempt)

		output, err = r.mainAction.Run(mCtx, output)
		if err == nil {
			return output, nil
		}

		if attempt < r.maxRetry {
			logger.Debugf(pCtx, "chain: rolling back with %s, error occurred: %v", r.rollbackAction.Name(), err)
			output, err = r.rollbackAction.Run(rCtx, output)
			if err != nil {
				return output, fmt.Errorf("rolling back failed: %w", err)
			}
		}
	}

	return output, err
}
