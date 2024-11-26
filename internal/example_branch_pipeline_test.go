package internal

import (
	"context"
	"github.com/JSYoo5B/chain"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBranchingPipeline(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)

	ctx := context.Background()
	t.Run("Basic Collatz with odd input(5)", func(t *testing.T) {
		simple := basicCollatzFunction()

		// 5 -odd-> (5*3=)15 -> (15+1=)16
		output, err := simple.Run(ctx, 5)

		assert.NoError(t, err)
		assert.Equal(t, 16, output)
	})
	t.Run("Basic Collatz with even input(16)", func(t *testing.T) {
		simple := basicCollatzFunction()

		// 16 -even-> (16/2=)8
		output, err := simple.Run(ctx, 16)

		assert.NoError(t, err)
		assert.Equal(t, 8, output)
	})
	t.Run("Shortcut Collatz with odd input(5)", func(t *testing.T) {
		shortcut := shortcutCollatzFunction()

		// 5 -odd-> (5*3=)15 -> (15+1=)16 -> (16/2=)8
		output, err := shortcut.Run(ctx, 5)

		assert.NoError(t, err)
		assert.Equal(t, 8, output)
	})
	t.Run("Shortcut Collatz with even input(16)", func(t *testing.T) {
		shortcut := shortcutCollatzFunction()

		// 16 -even-> (16/2=)8
		output, err := shortcut.Run(ctx, 16)

		assert.NoError(t, err)
		assert.Equal(t, 8, output)
	})
}

func basicCollatzFunction() *chain.Pipeline[int] {
	branch, even, odd1, odd2 := checkNext(), half(), triple(), inc()

	pipeline := chain.NewPipeline("SimpleCollatz", branch, even, odd1, odd2)
	pipeline.SetRunPlan(branch, chain.ActionPlan[int]{
		"even": even,
		"odd":  odd1,
	})
	pipeline.SetRunPlan(even, chain.TerminationPlan[int]())
	pipeline.SetRunPlan(odd1, chain.SuccessOnlyPlan(odd2))
	pipeline.SetRunPlan(odd2, chain.TerminationPlan[int]())

	return pipeline
}

func shortcutCollatzFunction() *chain.Pipeline[int] {
	branch, even, odd1, odd2 := checkNext(), half(), triple(), inc()

	pipeline := chain.NewPipeline("ShortcutCollatz", branch, even, odd1, odd2)
	pipeline.SetRunPlan(branch, chain.ActionPlan[int]{
		"even": even,
		"odd":  odd1,
	})
	pipeline.SetRunPlan(even, chain.TerminationPlan[int]())
	pipeline.SetRunPlan(odd1, chain.SuccessOnlyPlan(odd2))
	pipeline.SetRunPlan(odd2, chain.SuccessOnlyPlan(even))

	return pipeline
}

func checkNext() chain.Action[int] {
	branchFunc := func(_ context.Context, output int) (direction string, err error) {
		if output%2 == 0 {
			return "even", nil
		} else {
			return "odd", nil
		}
	}
	return chain.NewSimpleBranchAction("CheckNext", nil, []string{"even", "odd"}, branchFunc)
}

func half() chain.Action[int] {
	runFunc := func(_ context.Context, input int) (output int, err error) {
		return input / 2, nil
	}
	return chain.NewSimpleAction("Half", runFunc)
}

func triple() chain.Action[int] {
	runFunc := func(_ context.Context, input int) (output int, err error) {
		return input * 3, nil
	}
	return chain.NewSimpleAction("Triple", runFunc)
}

func inc() chain.Action[int] {
	runFunc := func(_ context.Context, input int) (output int, err error) {
		return input + 1, nil
	}
	return chain.NewSimpleAction("Inc", runFunc)
}
