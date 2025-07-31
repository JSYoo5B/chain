package logger

import (
	"context"
	"strings"
)

type runnerNameKey struct{}

func WithRunnerDepth(ctx context.Context, currentRunner string) context.Context {
	runnerName := currentRunner
	if parentName := ctx.Value(runnerNameKey{}); parentName != nil {
		runnerName = parentName.(string)
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
