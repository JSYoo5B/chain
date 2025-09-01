package chain

import (
	"context"
	"fmt"
	internalErrors "github.com/JSYoo5B/chain/internal/errors"
	"github.com/JSYoo5B/chain/internal/logger"
	"runtime/debug"
)

// AsRetryableAction creates an Action that retries the mainAction up to maxRetry times.
//
// If the mainAction fails, it executes the rollbackAction before the next retry attempt.
// The rollbackAction is not executed on the final attempt if it fails.
// rollbackAction can be skipped when it is nil.
func AsRetryableAction[T any](name string, mainAction, rollbackAction Action[T], maxRetry int) Action[T] {
	if mainAction == nil {
		panic("action cannot be nil")
	} else if maxRetry < 1 {
		panic("maxRetry must be greater than 0")
	}

	return &retryableAction[T]{
		name:           name,
		mainAction:     mainAction,
		rollbackAction: rollbackAction,
		maxRetry:       maxRetry,
	}
}

// SkipRollback provides an Action that explicitly skips the rollback process
// for a RetryableAction.
//
// When creating a RetryableAction with AsRetryableAction, use SkipRollback()
// as the rollback action to make the intent clear, rather than passing raw nil.
func SkipRollback[T any]() Action[T] {
	return nil
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

	var mCtx, rCtx context.Context
	mCtx = logger.WithRunnerDepth(pCtx, r.mainAction.Name())
	if r.rollbackAction != nil {
		rCtx = logger.WithRunnerDepth(pCtx, r.rollbackAction.Name())
	}
	for attempt := 1; attempt <= r.maxRetry; attempt++ {
		logger.Debugf(pCtx, "chain: running %s, attempt: %d", r.mainAction.Name(), attempt)

		output, err = r.mainAction.Run(mCtx, output)
		if err == nil {
			return output, nil
		}

		logger.Debugf(pCtx, "chain: error occurred: %v", err)

		if r.rollbackAction != nil && attempt < r.maxRetry {
			logger.Debugf(pCtx, "chain: rolling back with %s", r.rollbackAction.Name())
			output, err = r.rollbackAction.Run(rCtx, output)
			if err != nil {
				return output, fmt.Errorf("rolling back failed: %w", err)
			}
		}
	}

	return output, err
}
