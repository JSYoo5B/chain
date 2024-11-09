package dag

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"runtime/debug"
)

type Action[T any] interface {
	Name() string
	Run(ctx context.Context, input T) (output T, direction string, err error)
}

func TerminateAction[T any]() Action[T] {
	return nil
}

func runAction[T any](action Action[T], ctx context.Context, input T) (output T, direction string, runError error) {
	// Wrap panic handling for safe running in pipeline
	defer func() {
		if panicErr := recover(); panicErr != nil {
			logrus.Errorf("%s: panic occurred on running, caused by %s", action.Name(), panicErr)
			debug.PrintStack()

			output = input
			direction = Abort
			switch x := panicErr.(type) {
			case string:
				runError = errors.New(x)
			case error:
				runError = x
			default:
				runError = errors.New("unknown panic type")
			}
		}
	}()

	output, direction, runError = action.Run(ctx, input)
	if runError != nil && direction != Error {
		logrus.Errorf("%s: invoked error but directing '%s', overriding with Error", action.Name(), direction)
		direction = Error
	}
	return output, direction, runError
}
