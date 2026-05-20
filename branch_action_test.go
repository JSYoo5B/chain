package chain

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAsBranchAction(t *testing.T) {
	t.Run("wraps base action", func(t *testing.T) {
		baseAction := NewSimpleAction(
			"double",
			func(_ context.Context, input int) (int, error) {
				return input * 2, nil
			},
		)
		action := AsBranchAction(
			baseAction,
			func(_ context.Context, output int) (string, error) {
				if output%4 == 0 {
					return "evenlyDivided", nil
				}
				return "remaining", nil
			},
			"evenlyDivided",
			"remaining",
		)

		output, err := action.Run(context.Background(), 2)
		direction, directionErr := action.NextDirection(context.Background(), output)

		assert.NoError(t, err)
		assert.NoError(t, directionErr)
		assert.Equal(t, "double", action.Name())
		assert.Equal(t, 4, output)
		assert.Equal(t, "evenlyDivided", direction)
		assert.Equal(t, []string{"evenlyDivided", "remaining"}, action.Directions())
	})

	t.Run("directions returns copy", func(t *testing.T) {
		baseAction := NewSimpleAction(
			"pass",
			func(_ context.Context, input int) (int, error) {
				return input, nil
			},
		)
		action := AsBranchAction(
			baseAction,
			func(_ context.Context, _ int) (string, error) {
				return "left", nil
			},
			"left",
			"right",
		)

		directions := action.Directions()
		directions[0] = "changed"

		assert.Equal(t, []string{"left", "right"}, action.Directions())
	})

	t.Run("allows built-in directions without custom directions", func(t *testing.T) {
		baseAction := NewSimpleAction(
			"pass",
			func(_ context.Context, input int) (int, error) {
				return input, nil
			},
		)
		action := AsBranchAction(
			baseAction,
			func(_ context.Context, _ int) (string, error) {
				return Success, nil
			},
		)

		assert.Empty(t, action.Directions())
	})
}
