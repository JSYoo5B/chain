package chain

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRecoverOnRun(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)

	divide := func(name string) Action[int] {
		return NewSimpleAction(
			name,
			func(_ context.Context, input int) (int, error) { return 1 / input, nil },
		)
	}
	setTen := func(name string) Action[int] {
		return NewSimpleAction(
			name,
			func(_ context.Context, _ int) (int, error) { return 10, nil },
		)
	}

	ctx := context.Background()
	type testCase struct {
		closure func(t *testing.T) func()
	}

	testCases := map[string]testCase{
		"recover from action": {
			func(t *testing.T) func() {
				return func() {
					pipeline := NewPipeline("Calculation", divide("1"))
					_, err := pipeline.Run(ctx, 0)

					assert.Error(t, err)
					assert.Contains(t, err.Error(), "divide by zero")
				}
			}},
		"skip actions after recover": {
			func(t *testing.T) func() {
				return func() {
					pipeline := NewPipeline("Calculation", divide("1"), setTen("2"))
					output, err := pipeline.Run(ctx, 0)

					assert.Error(t, err)
					assert.Contains(t, err.Error(), "divide by zero")
					assert.NotEqual(t, 10, output)
				}
			}},
		"internal pipeline panic recovers": {
			func(t *testing.T) func() {
				return func() {
					subPipeline := NewPipeline("SubPipeline", divide(".1"), setTen(".2"))
					pipeline := NewPipeline(
						"SuperPipeline",
						setTen("1"),
						divide("2"),
						subPipeline,
						setTen("4"),
					)
					output, err := pipeline.Run(ctx, 0)

					assert.Error(t, err)
					assert.Contains(t, err.Error(), "divide by zero")
					assert.NotEqual(t, 10, output)
				}
			}},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			code := tc.closure(t)
			assert.NotPanics(t, code)
		})
	}
}
