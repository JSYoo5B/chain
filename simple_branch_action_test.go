package chain

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleBranchAction_Directions(t *testing.T) {
	branchFunc := func(_ context.Context, _ int) (string, error) {
		return "left", nil
	}

	t.Run("constructor copies directions", func(t *testing.T) {
		directions := []string{"left", "right"}
		action := NewSimpleBranchAction("branch", nil, directions, branchFunc)

		directions[0] = "changed"

		assert.Equal(t, []string{"left", "right"}, action.Directions())
	})

	t.Run("directions returns copy", func(t *testing.T) {
		action := NewSimpleBranchAction("branch", nil, []string{"left", "right"}, branchFunc)
		directions := action.Directions()

		directions[0] = "changed"

		assert.Equal(t, []string{"left", "right"}, action.Directions())
	})
}
