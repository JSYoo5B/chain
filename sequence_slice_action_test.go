package chain

import (
	"context"
	"errors"
	"github.com/JSYoo5B/chain/internal/logger"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestSequenceSliceAction(t *testing.T) {
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
		doubles := AsSequenceSliceAction("MapDouble", double, false)
		input := []int{1, 2, 3, 4, 5}
		expected := []int{2, 4, 6, 8, 10}

		output, err := doubles.Run(context.Background(), input)

		assert.NoError(t, err)
		assert.Equal(t, expected, output)
	})
	t.Run("error in iteration stops", func(t *testing.T) {
		doubles := AsSequenceSliceAction("MapDoubleStop", positiveDouble, true)
		input := []int{1, 2, -1, 4, 5}
		expected := []int{2, 4, 0, 4, 5}

		output, err := doubles.Run(context.Background(), input)

		assert.Error(t, err)
		assert.Equal(t, expected, output)
	})
	t.Run("error in iteration continues", func(t *testing.T) {
		doubles := AsSequenceSliceAction("MapDoubleContinue", positiveDouble, false)
		input := []int{1, 2, -1, 4, 5}
		expected := []int{2, 4, 0, 8, 10}

		output, err := doubles.Run(context.Background(), input)

		assert.Error(t, err)
		assert.Equal(t, expected, output)
	})
	t.Run("multiple errors in iteration are joined in order", func(t *testing.T) {
		onlyPositive := NewSimpleAction(
			"onlyPositive",
			func(_ context.Context, input int) (int, error) {
				if input < 0 {
					return 0, errors.New("negative input")
				}
				return input, nil
			})
		action := AsSequenceSliceAction("Sequence", onlyPositive, false)
		input := []int{1, -1, 2, -2, 3}
		expected := []int{1, 0, 2, 0, 3}

		output, err := action.Run(context.Background(), input)

		assert.Error(t, err)
		assert.Equal(t, expected, output)
		assert.Contains(t, err.Error(), "negative input")
		assert.Less(t, strings.Index(err.Error(), "index 1"), strings.Index(err.Error(), "index 3"))
	})
	t.Run("error before panic is preserved", func(t *testing.T) {
		failOrDivide := NewSimpleAction(
			"failOrDivide",
			func(_ context.Context, input int) (int, error) {
				if input < 0 {
					return 0, errors.New("negative input")
				}
				return 10 / input, nil
			})
		action := AsSequenceSliceAction("Sequence", failOrDivide, false)
		input := []int{10, -1, 5, 0, 1}
		expected := []int{1, 0, 2, 0, 1}

		output, err := action.Run(context.Background(), input)

		assert.Error(t, err)
		assert.Equal(t, expected, output)
		assert.Contains(t, err.Error(), "negative input")
		assert.Contains(t, err.Error(), "divide by zero")
		assert.Less(t, strings.Index(err.Error(), "negative input"), strings.Index(err.Error(), "divide by zero"))
	})
	t.Run("panic in iteration", func(t *testing.T) {
		divides := AsSequenceSliceAction("MapDivide10", divide10, false)
		input := []int{10, 5, 2, 0, 1}
		expected := []int{1, 2, 5, 0, 1}

		output, err := divides.Run(context.Background(), input)

		assert.Error(t, err)
		assert.Equal(t, expected, output)
	})
}
