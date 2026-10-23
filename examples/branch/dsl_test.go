package branch

import (
	"context"
	"testing"

	"github.com/JSYoo5B/chain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGeneratedCollatzMatchesManual(t *testing.T) {
	for name, builders := range map[string][2]func() *chain.Workflow[int]{
		"basic":    {basicCollatzFunction, generatedBasicCollatzFunction},
		"shortcut": {shortcutCollatzFunction, generatedShortcutCollatzFunction},
	} {
		t.Run(name, func(t *testing.T) {
			manual, generated := builders[0](), builders[1]()
			require.NoError(t, generated.ValidateGraph())
			assert.Equal(t, manual.Name(), generated.Name())
			for _, input := range []int{-5, -2, 0, 1, 5, 16, 27} {
				want, wantErr := manual.Run(context.Background(), input)
				got, gotErr := generated.Run(context.Background(), input)
				require.NoError(t, wantErr)
				require.NoError(t, gotErr)
				assert.Equal(t, want, got, "input %d", input)
			}
		})
	}
}

func TestGeneratedLeafDoesNotFallThrough(t *testing.T) {
	action := func(name string, value int) chain.Action[int] {
		return chain.NewSimpleAction(name, func(_ context.Context, input int) (int, error) {
			return input + value, nil
		})
	}
	w := generatedLeaf(action("first", 1), action("leaf", 10), action("alternate", 100))
	require.NoError(t, w.ValidateGraph())
	got, err := w.Run(context.Background(), 0)
	require.NoError(t, err)
	assert.Equal(t, 11, got)
}

func TestGeneratedPrerequisiteClosure(t *testing.T) {
	w := generatedIncrement()
	require.NoError(t, w.ValidateGraph())
	got, err := w.Run(context.Background(), 41)
	require.NoError(t, err)
	assert.Equal(t, 42, got)
}
