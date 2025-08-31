package chain

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRetryableAction(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)

	checkZero := NewSimpleAction(
		"checkZero",
		func(ctx context.Context, input int) (int, error) {
			if input == 0 {
				return input, nil
			}
			return input, fmt.Errorf("%d is not zero", input)
		})
	decrease := NewSimpleAction(
		"decrase",
		func(ctx context.Context, input int) (int, error) {
			return input - 1, nil
		})
	expectZero := AsRetryableAction("expectZero", checkZero, decrease, 3)

	t.Run("direct success", func(t *testing.T) {
		output, err := expectZero.Run(context.Background(), 0)

		assert.NoError(t, err)
		assert.Equal(t, 0, output)
	})

	t.Run("first retry success", func(t *testing.T) {
		output, err := expectZero.Run(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, 0, output)
	})

	t.Run("max retry success", func(t *testing.T) {
		output, err := expectZero.Run(context.Background(), 2)

		assert.NoError(t, err)
		assert.Equal(t, 0, output)
	})

	t.Run("max retry fail", func(t *testing.T) {
		output, err := expectZero.Run(context.Background(), 10)

		assert.Error(t, err)
		assert.NotEqual(t, 0, output)
	})
}

func TestRetryableAction_withoutRollback(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)

	checkZeroAndDecrease := NewSimpleAction(
		"checkZeroAndDecrease",
		func(ctx context.Context, input int) (int, error) {
			if input == 0 {
				return input, nil
			}
			return input - 1, fmt.Errorf("%d was not zero", input)
		})
	expectZero := AsRetryableAction("expectZero", checkZeroAndDecrease, nil, 3)

	t.Run("direct success", func(t *testing.T) {
		output, err := expectZero.Run(context.Background(), 0)

		assert.NoError(t, err)
		assert.Equal(t, 0, output)
	})

	t.Run("first retry success", func(t *testing.T) {
		output, err := expectZero.Run(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, 0, output)
	})

	t.Run("max retry success", func(t *testing.T) {
		output, err := expectZero.Run(context.Background(), 2)

		assert.NoError(t, err)
		assert.Equal(t, 0, output)
	})

	t.Run("max retry fail", func(t *testing.T) {
		output, err := expectZero.Run(context.Background(), 10)

		assert.Error(t, err)
		assert.NotEqual(t, 0, output)
	})
}
