package chain

import (
	"context"
	"errors"
	"github.com/JSYoo5B/chain/internal/logger"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParallelSlicePipeline(t *testing.T) {
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

	t.Run("simple iteration", func(t *testing.T) {
		doubles := NewParallelSlicePipeline("MapDouble", double)
		input := []int{1, 2, 3, 4, 5}
		expected := []int{2, 4, 6, 8, 10}

		output, err := doubles.Run(context.Background(), input)

		assert.NoError(t, err)
		assert.Equal(t, expected, output)
	})
	t.Run("error in iteration", func(t *testing.T) {
		doubles := NewParallelSlicePipeline("MapDoubleStop", positiveDouble)
		input := []int{1, 2, -1, 4, 5}
		expected := []int{2, 4, 0, 8, 10}

		output, err := doubles.Run(context.Background(), input)

		assert.Error(t, err)
		assert.Equal(t, expected, output)
	})
	t.Run("panic in iteration", func(t *testing.T) {
		divides := NewParallelSlicePipeline("MapDivide10", divide10)
		input := []int{10, 5, 2, 0, 1}
		expected := []int{1, 2, 5, 0, 10}

		output, err := divides.Run(context.Background(), input)

		assert.Error(t, err)
		assert.Equal(t, expected, output)
	})
}
