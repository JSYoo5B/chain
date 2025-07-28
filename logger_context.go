package chain

import (
	"context"
	"strings"
)

type runnerNameKey struct{}

func appendRunnerName(ctx context.Context, currentRunner string) context.Context {
	runnerName := currentRunner
	if parentName := ctx.Value(runnerNameKey{}); parentName != nil {
		runnerName = parentName.(string)
	}

	if !strings.HasSuffix(runnerName, currentRunner) {
		runnerName += "/" + currentRunner
	}

	return context.WithValue(ctx, runnerNameKey{}, runnerName)
}
