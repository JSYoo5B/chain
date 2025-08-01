package logger

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWithRunnerDepth(t *testing.T) {
	t.Run("nil context should panic", func(t *testing.T) {
		assert.Panics(t, func() {
			WithRunnerDepth(nil, "test")
		})
	})

	t.Run("without parent runner", func(t *testing.T) {
		ctx := WithRunnerDepth(context.Background(), "test")

		runnerName, ok := RunnerNameFromContext(ctx)

		assert.True(t, ok)
		assert.Equal(t, "test", runnerName)
	})

	t.Run("parent runner", func(t *testing.T) {
		parentCtx := WithRunnerDepth(context.Background(), "parent")
		childCtx := WithRunnerDepth(parentCtx, "child")

		runnerName, ok := RunnerNameFromContext(childCtx)

		assert.True(t, ok)
		assert.Equal(t, "parent/child", runnerName)
	})

	t.Run("duplicate parent runner", func(t *testing.T) {
		parentCtx := WithRunnerDepth(context.Background(), "parent")
		childCtx := WithRunnerDepth(parentCtx, "child")
		duplicateCtx := WithRunnerDepth(childCtx, "child")

		runnerName, ok := RunnerNameFromContext(duplicateCtx)

		assert.True(t, ok)
		assert.Equal(t, "parent/child", runnerName)
	})
}

func TestRunnerNameFromContext(t *testing.T) {
	t.Run("nil context", func(t *testing.T) {
		_, ok := RunnerNameFromContext(nil)
		assert.False(t, ok)
	})
	t.Run("empty context", func(t *testing.T) {
		_, ok := RunnerNameFromContext(context.Background())
		assert.False(t, ok)
	})
	t.Run("context with simple runner name", func(t *testing.T) {
		ctx := WithRunnerDepth(context.Background(), "test")
		runnerName, ok := RunnerNameFromContext(ctx)
		assert.True(t, ok)
		assert.Equal(t, "test", runnerName)
	})
	t.Run("context with complex runner name", func(t *testing.T) {
		ctx := WithRunnerDepth(context.Background(), "test/inner")
		runnerName, ok := RunnerNameFromContext(ctx)
		assert.True(t, ok)
		assert.Equal(t, "test/inner", runnerName)
	})
}
