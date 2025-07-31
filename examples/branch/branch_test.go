package branch

import (
	"context"
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
