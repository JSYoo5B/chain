package errorcontrol

import (
	"context"
	"errors"

	"github.com/JSYoo5B/chain"
)

type attemptCounts struct {
	Calls     int
	Rollbacks int
}

// The first two additions fail; the third succeeds. Rollback can undo each
// failed addition before the next attempt.
func retryingIncrement(useRollback bool) (chain.Action[int], *attemptCounts) {
	counts := &attemptCounts{}
	step := chain.NewSimpleAction("add one", func(_ context.Context, input int) (int, error) {
		counts.Calls++
		if counts.Calls < 3 {
			return input + 1, errors.New("not yet")
		}
		return input + 1, nil
	})
	undo := chain.SkipRollback[int]()
	if useRollback {
		undo = chain.NewSimpleAction("subtract one", func(_ context.Context, input int) (int, error) {
			counts.Rollbacks++
			return input - 1, nil
		})
	}
	return chain.AsRetryableAction("retry", step, undo, 3), counts
}

type fallbackRecord struct {
	Input int
	Err   error
}

func optionalIncrement() (chain.Action[int], *fallbackRecord) {
	record := &fallbackRecord{}
	step := chain.NewSimpleAction("add one", func(_ context.Context, input int) (int, error) {
		return input + 1, errors.New("optional step failed")
	})
	optional := chain.AsBestEffortAction(step, func(_ context.Context, input int, err error) {
		record.Input, record.Err = input, err
	})
	return optional, record
}
