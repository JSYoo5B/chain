package chain

import "context"

// BestEffortFallbackFunc is called when a BestEffortAction's wrapped action returns an error.
// The input argument is the value that was passed to the wrapped action.
type BestEffortFallbackFunc[T any] func(ctx context.Context, input T, err error)

// AsBestEffortAction creates an Action that suppresses errors from the wrapped action.
//
// The wrapped action keeps its original Name behavior. When the wrapped action succeeds,
// its output is returned normally. When it returns an error, fallback is called with the
// original input and error, then the wrapped action's output is returned with a nil error.
// This allows a Workflow to continue through the Success direction even when the wrapped
// action failed in a non-critical way.
func AsBestEffortAction[T any](action Action[T], fallback BestEffortFallbackFunc[T]) Action[T] {
	return &bestEffortAction[T]{
		Action:   action,
		fallback: fallback,
	}
}

type bestEffortAction[T any] struct {
	Action[T]
	fallback BestEffortFallbackFunc[T]
}

func (b bestEffortAction[T]) Run(ctx context.Context, input T) (output T, err error) {
	output, err = b.Action.Run(ctx, input)
	if err == nil {
		return output, nil
	}

	if b.fallback != nil {
		b.fallback(ctx, input, err)
	}
	return output, nil
}
