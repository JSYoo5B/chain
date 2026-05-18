package chain

import (
	"context"
	"fmt"
	"github.com/JSYoo5B/chain/internal/errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"strings"
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

func TestRetryableAction_RollbackError(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)

	checkZero := NewSimpleAction(
		"checkZero",
		func(ctx context.Context, input int) (int, error) {
			if input == 0 {
				return input, nil
			}
			return input, fmt.Errorf("%d is not zero", input)
		})

	t.Run("rollback error preserves main action error", func(t *testing.T) {
		rollbackFails := NewSimpleAction(
			"rollbackFails",
			func(_ context.Context, input int) (int, error) {
				return input - 1, fmt.Errorf("rollback failed")
			})
		expectZero := AsRetryableAction("expectZero", checkZero, rollbackFails, 3)

		output, err := expectZero.Run(context.Background(), 2)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "2 is not zero")
		assert.Contains(t, err.Error(), "rolling back failed")
		assert.Contains(t, err.Error(), "rollback failed")
		assertLessIndex(t, err.Error(), "2 is not zero", "rolling back failed")
		assert.Equal(t, 1, output)
	})

	t.Run("rollback panic preserves main action error", func(t *testing.T) {
		rollbackPanics := NewSimpleAction(
			"rollbackPanics",
			func(_ context.Context, input int) (int, error) {
				panic("rollback panic")
			})
		expectZero := AsRetryableAction("expectZero", checkZero, rollbackPanics, 3)

		output, err := expectZero.Run(context.Background(), 2)

		var panicErr *errors.PanicError
		assert.Error(t, err)
		assert.ErrorAs(t, err, &panicErr)
		assert.Contains(t, err.Error(), "2 is not zero")
		assert.Contains(t, err.Error(), "rollback panic")
		assertLessIndex(t, err.Error(), "2 is not zero", "rollback panic")
		assert.Equal(t, 2, output)
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
	expectZero := AsRetryableAction("expectZero", checkZeroAndDecrease, SkipRollback[int](), 3)

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

func assertLessIndex(t *testing.T, s, before, after string) {
	t.Helper()

	beforeIndex := strings.Index(s, before)
	afterIndex := strings.Index(s, after)
	if assert.NotEqual(t, -1, beforeIndex) && assert.NotEqual(t, -1, afterIndex) {
		assert.Less(t, beforeIndex, afterIndex)
	}
}
