package internal

import (
	"context"
	"fmt"
	"github.com/JSYoo5B/dago/railway"
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
		output, direction, err := simple.Run(ctx, 5)

		assert.NoError(t, err)
		assert.Equal(t, 16, output)
		assert.Equal(t, railway.Success, direction)
	})
	t.Run("Basic Collatz with even input(16)", func(t *testing.T) {
		simple := basicCollatzFunction()

		// 16 -even-> (16/2=)8
		output, direction, err := simple.Run(ctx, 16)

		assert.NoError(t, err)
		assert.Equal(t, 8, output)
		assert.Equal(t, railway.Success, direction)
	})
	t.Run("Shortcut Collatz with odd input(5)", func(t *testing.T) {
		shortcut := shortcutCollatzFunction()

		// 5 -odd-> (5*3=)15 -> (15+1=)16 -> (16/2=)8
		output, direction, err := shortcut.Run(ctx, 5)

		assert.NoError(t, err)
		assert.Equal(t, 8, output)
		assert.Equal(t, railway.Success, direction)
	})
	t.Run("Shortcut Collatz with even input(16)", func(t *testing.T) {
		shortcut := shortcutCollatzFunction()

		// 16 -even-> (16/2=)8
		output, direction, err := shortcut.Run(ctx, 16)

		assert.NoError(t, err)
		assert.Equal(t, 8, output)
		assert.Equal(t, railway.Success, direction)
	})
}

func basicCollatzFunction() *railway.Pipeline[int] {
	branch, even, odd1, odd2 := checkNext(), half(), triple(), inc()

	pipeline := railway.NewPipeline("SimpleCollatz", branch, even, odd1, odd2)
	pipeline.SetRunPlan(branch, railway.ActionPlan[int]{
		"even": even,
		"odd":  odd1,
	})
	pipeline.SetRunPlan(even, railway.TerminationPlan[int]())
	pipeline.SetRunPlan(odd1, railway.SuccessOnlyPlan(odd2))
	pipeline.SetRunPlan(odd2, railway.TerminationPlan[int]())

	return pipeline
}

func shortcutCollatzFunction() *railway.Pipeline[int] {
	branch, even, odd1, odd2 := checkNext(), half(), triple(), inc()

	pipeline := railway.NewPipeline("ShortcutCollatz", branch, even, odd1, odd2)
	pipeline.SetRunPlan(branch, railway.ActionPlan[int]{
		"even": even,
		"odd":  odd1,
	})
	pipeline.SetRunPlan(even, railway.TerminationPlan[int]())
	pipeline.SetRunPlan(odd1, railway.SuccessOnlyPlan(odd2))
	pipeline.SetRunPlan(odd2, railway.SuccessOnlyPlan(even))

	return pipeline
}

func checkNext() railway.Action[int] {
	runFunc := func(_ context.Context, input int) (output int, direction string, err error) {
		if input%2 == 0 {
			return input, "even", nil
		} else {
			return input, "odd", nil
		}
	}
	return railway.NewSimpleBranchAction("CheckNext", []string{"even", "odd"}, runFunc)
}

func half() railway.Action[int] {
	runFunc := func(_ context.Context, input int) (output int, direction string, err error) {
		if input%2 != 0 {
			return input, railway.Error, fmt.Errorf("cannot try half with odd")
		}
		return input / 2, railway.Success, nil
	}
	return railway.NewSimpleAction("Half", runFunc)
}

func triple() railway.Action[int] {
	runFunc := func(_ context.Context, input int) (output int, direction string, err error) {
		return input * 3, railway.Success, nil
	}
	return railway.NewSimpleAction("Triple", runFunc)
}

func inc() railway.Action[int] {
	runFunc := func(_ context.Context, input int) (output int, direction string, err error) {
		return input + 1, railway.Success, nil
	}
	return railway.NewSimpleAction("Inc", runFunc)
}
