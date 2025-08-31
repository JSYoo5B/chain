package logger

import (
	"context"
	"strings"
)

type runnerNameKey struct{}

func WithRunnerDepth(ctx context.Context, currentRunner string) context.Context {
	if ctx == nil {
		panic("cannot create context from nil parent")
	} else if currentRunner == "" {
		return ctx
	}

	runnerName := currentRunner
	if parentName, ok := ctx.Value(runnerNameKey{}).(string); ok {
		runnerName = parentName
	}

	if !strings.HasSuffix(runnerName, currentRunner) {
		runnerName += "/" + currentRunner
	}

	return context.WithValue(ctx, runnerNameKey{}, runnerName)
}

func RunnerNameFromContext(ctx context.Context) (string, bool) {
	if ctx == nil {
		return "", false
	}

	runnerName, ok := ctx.Value(runnerNameKey{}).(string)
	return runnerName, ok
}
