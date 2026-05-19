package chain

import (
	"context"
	"errors"
	"fmt"
	"github.com/JSYoo5B/chain/internal/logger"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParallelSliceAction(t *testing.T) {
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
		doubles := AsParallelSliceAction("MapDouble", double)
		input := []int{1, 2, 3, 4, 5}
		expected := []int{2, 4, 6, 8, 10}

		output, err := doubles.Run(context.Background(), input)

		assert.NoError(t, err)
		assert.Equal(t, expected, output)
	})
	t.Run("error in parallel iteration continues", func(t *testing.T) {
		doubles := AsParallelSliceAction("MapDoubleContinue", positiveDouble)
		input := []int{1, 2, -1, 4, 5}
		expected := []int{2, 4, 0, 8, 10}

		output, err := doubles.Run(context.Background(), input)

		assert.Error(t, err)
		assert.Equal(t, expected, output)
	})
	t.Run("multiple errors are joined", func(t *testing.T) {
		failSmallNumbers := NewSimpleAction(
			"failSmallNumbers",
			func(_ context.Context, input int) (int, error) {
				if input < 2 {
					return input, fmt.Errorf("bad input: %d", input)
				}
				return input * 2, nil
			})
		doubles := AsParallelSliceAction("MapDoubleWithErrors", failSmallNumbers)
		input := []int{0, 1, 2, 3}
		expected := []int{0, 1, 4, 6}

		output, err := doubles.Run(context.Background(), input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "bad input: 0")
		assert.Contains(t, err.Error(), "bad input: 1")
		assert.Equal(t, expected, output)
	})
	t.Run("panic and error are joined", func(t *testing.T) {
		failAndPanic := NewSimpleAction(
			"failAndPanic",
			func(_ context.Context, input int) (int, error) {
				switch input {
				case 0:
					return input, fmt.Errorf("bad input: %d", input)
				case 1:
					panic("panic input: 1")
				default:
					return input * 2, nil
				}
			})
		doubles := AsParallelSliceAction("MapDoubleWithPanicAndError", failAndPanic)
		input := []int{0, 1, 2}
		expected := []int{0, 1, 4}

		output, err := doubles.Run(context.Background(), input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "bad input: 0")
		assert.Contains(t, err.Error(), "panic input: 1")
		assert.Equal(t, expected, output)
	})
	t.Run("panic in parallel iteration", func(t *testing.T) {
		divides := AsParallelSliceAction("MapDivide10", divide10)
		input := []int{10, 5, 2, 0, 1}
		expected := []int{1, 2, 5, 0, 10}

		output, err := divides.Run(context.Background(), input)

		assert.Error(t, err)
		assert.Equal(t, expected, output)
	})
}
