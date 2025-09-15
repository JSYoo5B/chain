package chain

import (
	"context"
	"errors"
	"fmt"
	internalErrors "github.com/JSYoo5B/chain/internal/errors"
	"github.com/JSYoo5B/chain/internal/logger"
	"runtime/debug"
)

// Run executes the Workflow by running Actions in the order they were configured,
// starting from the initAction, which is the first one of the memberActions provided
// by the constructor such as NewWorkflow.
// The actions are executed in order, passing the output of one action as input to the next.
func (w *Workflow[T]) Run(ctx context.Context, input T) (output T, err error) {
	if len(w.runPlans) == 1 {
		output, _, err = runAction(w.initAction, ctx, input)
		return output, err
	}

	return w.RunAt(w.initAction, ctx, input)
}

// RunAt starts the execution of the Workflow from a given Action (initAction).
// It follows the action plan, executing actions sequentially based on the specified directions.
// If an action returns an error, the Workflow will proceed to the next action according to
// the defined plan, potentially directing the flow to an action mapped for the Failure direction.
// The Abort direction, when encountered, will immediately halt the Workflow execution unless
// the plan specifies otherwise.
// If no action plan is found for a given direction,
// the Workflow will terminate with the appropriate error.
func (w *Workflow[T]) RunAt(initAction Action[T], ctx context.Context, input T) (output T, lastErr error) {
	if !isMemberActionInWorkflow(initAction, w) {
		return input, errors.New("given initAction is not registered on constructor")
	}

	ctx = logger.WithRunnerDepth(ctx, w.name)

	var (
		terminate     = Terminate[T]()
		currentAction Action[T]
		nextAction    Action[T]
		direction     string
		runErr        error
		selectErr     error
	)
	logger.Debugf(ctx, "chain: start running with `%s`", initAction.Name())
	for currentAction = initAction; currentAction != nil; currentAction = nextAction {
		output, direction, runErr = runAction(currentAction, ctx, input)

		nextAction, selectErr = selectNextAction(w.runPlans[currentAction], currentAction, direction)
		if selectErr != nil {
			logger.Error(ctx, selectErr)
			direction = Abort
			lastErr = selectErr
			break
		}

		nextActionName := "termination"
		if nextAction != terminate {
			nextActionName = nextAction.Name()
		}
		logger.Debugf(ctx, "chain: `%s` directs `%s`, selecting `%s`", currentAction.Name(), direction, nextActionName)

		input = output
		if runErr != nil {
			lastErr = runErr
		}
	}
	if lastErr != nil && direction != Abort {
		direction = Failure
	}

	return output, lastErr
}

func selectNextAction[T any](plan RunPlan[T], currentAction Action[T], direction string) (nextAction Action[T], err error) {
	var (
		terminate = Terminate[T]()
		exist     bool
	)
	if plan == nil {
		return terminate, fmt.Errorf("no action plan found for `%s`", currentAction.Name())
	}
	if nextAction, exist = plan[direction]; !exist {
		return terminate, fmt.Errorf("no action plan from `%s` directing `%s`", currentAction.Name(), direction)
	}

	return nextAction, nil
}

func runAction[T any](action Action[T], ctx context.Context, input T) (output T, direction string, runError error) {
	ctx = logger.WithRunnerDepth(ctx, action.Name())
	runnerName, _ := logger.RunnerNameFromContext(ctx)

	// Wrap panic handling for safe running in a Workflow
	defer func() {
		if panicErr := recover(); panicErr != nil {
			logger.Errorf(ctx, "chain: panic occurred on running, caused by %v", panicErr)
			debug.PrintStack()
			output, direction = input, Abort
			runError = internalErrors.NewPanicError(runnerName, panicErr)
		}
	}()

	output, runError = action.Run(ctx, input)
	if runError != nil {
		var panicError *internalErrors.PanicError
		if errors.As(runError, &panicError) {
			return output, Abort, runError
		}

		return output, Failure, runError
	}
	direction = Success
	if branchAction, isBranchAction := action.(BranchAction[T]); isBranchAction {
		direction, runError = branchAction.NextDirection(ctx, output)
	}

	return output, direction, runError
}
