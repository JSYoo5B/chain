package chain

import (
	"context"
	"errors"
	"github.com/JSYoo5B/chain/internal/logger"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParallelMapAction(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)

	double := NewSimpleAction(
		"double",
		func(ctx context.Context, input int) (int, error) {
			logger.Debugf(ctx, "doubling %d", input)
			return input * 2, nil
		})
	positiveDouble := NewSimpleAction(
		"positiveDouble",
		func(ctx context.Context, input int) (int, error) {
			logger.Debugf(ctx, "doubling %d (only positives)", input)
			if input < 0 {
				return 0, errors.New("negative input")
			}
			return input * 2, nil
		})
	divide10 := NewSimpleAction(
		"divide10",
		func(ctx context.Context, input int) (int, error) {
			logger.Debugf(ctx, "dividing 10 with %d", input)
			return 10 / input, nil
		})

	t.Run("simple parallel iteration", func(t *testing.T) {
		doubles := NewParallelMapAction[string, int]("MapDouble", double)
		input := map[string]int{"one": 1, "two": 2, "three": 3, "four": 4, "five": 5}
		expected := map[string]int{"one": 2, "two": 4, "three": 6, "four": 8, "five": 10}

		output, err := doubles.Run(context.Background(), input)

		assert.NoError(t, err)
		assert.Equal(t, expected, output)
	})
	t.Run("error in parallel iteration continues", func(t *testing.T) {
		doubles := NewParallelMapAction[string, int]("MapDoubleContinue", positiveDouble)
		input := map[string]int{"one": 1, "two": 2, "minus": -1, "four": 4, "five": 5}
		expected := map[string]int{"one": 2, "two": 4, "minus": 0, "four": 8, "five": 10}

		output, err := doubles.Run(context.Background(), input)

		assert.Error(t, err)
		assert.Equal(t, expected, output)
	})
	t.Run("panic in parallel iteration", func(t *testing.T) {
		divides := NewParallelMapAction[string, int]("MapDivide10", divide10)
		input := map[string]int{"ten": 10, "five": 5, "two": 2, "zero": 0, "one": 1}
		expected := map[string]int{"ten": 1, "five": 2, "two": 5, "zero": 0, "one": 10}

		output, err := divides.Run(context.Background(), input)

		assert.Error(t, err)
		assert.Equal(t, expected, output)
	})
}
