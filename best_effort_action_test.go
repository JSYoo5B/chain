package chain

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBestEffortAction(t *testing.T) {
	t.Run("returns wrapped action output on success", func(t *testing.T) {
		action := NewSimpleAction(
			"double",
			func(_ context.Context, input int) (int, error) {
				return input * 2, nil
			},
		)
		bestEffort := AsBestEffortAction(action, nil)

		output, err := bestEffort.Run(context.Background(), 2)

		assert.NoError(t, err)
		assert.Equal(t, "double", bestEffort.Name())
		assert.Equal(t, 4, output)
	})

	t.Run("suppresses wrapped action error and calls fallback", func(t *testing.T) {
		actionErr := errors.New("non-critical failure")
		action := NewSimpleAction(
			"fail",
			func(_ context.Context, input int) (int, error) {
				return input + 1, actionErr
			},
		)
		var fallbackInput int
		var fallbackErr error
		bestEffort := AsBestEffortAction(
			action,
			func(_ context.Context, input int, err error) {
				fallbackInput = input
				fallbackErr = err
			},
		)

		output, err := bestEffort.Run(context.Background(), 2)

		assert.NoError(t, err)
		assert.Equal(t, 3, output)
		assert.Equal(t, 2, fallbackInput)
		assert.ErrorIs(t, fallbackErr, actionErr)
	})

	t.Run("continues workflow through success direction after wrapped action error", func(t *testing.T) {
		actionErr := errors.New("non-critical failure")
		fail := AsBestEffortAction(
			NewSimpleAction(
				"fail",
				func(_ context.Context, input int) (int, error) {
					return input + 1, actionErr
				},
			),
			nil,
		)
		success := NewSimpleAction(
			"success",
			func(_ context.Context, input int) (int, error) {
				return input * 2, nil
			},
		)
		failure := NewSimpleAction(
			"failure",
			func(_ context.Context, input int) (int, error) {
				return 0, errors.New("should not run")
			},
		)
		workflow := NewWorkflow("workflow", fail, success, failure)
		workflow.SetRunPlan(fail, DefaultPlan(success, failure))
		workflow.SetRunPlan(success, TerminationPlan[int]())

		output, err := workflow.Run(context.Background(), 2)

		assert.NoError(t, err)
		assert.Equal(t, 6, output)
	})
}
