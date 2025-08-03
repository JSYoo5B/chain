package chain

import (
	"context"
	"errors"
	"fmt"
	"github.com/JSYoo5B/chain/internal/logger"
	"runtime/debug"
)

// Run executes the Pipeline by running Actions in the order they were configured,
// starting from the initAction, which is the first one of the memberActions provided
// by the constructor such as NewPipeline.
// The actions are executed in order, passing the output of one action as input to the next.
func (p *Pipeline[T]) Run(ctx context.Context, input T) (output T, err error) {
	if len(p.runPlans) == 1 {
		output, _, err = runAction(p.initAction, ctx, input)
		return output, err
	}

	return p.RunAt(p.initAction, ctx, input)
}

// RunAt starts the execution of the pipeline from a given Action (initAction).
// It follows the action plan, executing actions sequentially based on the specified directions.
// If an action returns an error, the pipeline will proceed to the next action according to
// the defined plan, potentially directing the flow to an action mapped for the Error direction.
// The Abort direction, when encountered, will immediately halt the pipeline execution unless
// the plan specifies otherwise.
// If no action plan is found for a given direction,
// the pipeline will terminate with the appropriate error.
func (p *Pipeline[T]) RunAt(initAction Action[T], ctx context.Context, input T) (output T, lastErr error) {
	if !isMemberActionInPipeline(initAction, p) {
		return input, errors.New("given initAction is not registered on constructor")
	}

	ctx = logger.WithRunnerDepth(ctx, p.name)

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

		nextAction, selectErr = selectNextAction(p.runPlans[currentAction], currentAction, direction)
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
		direction = Error
	}

	return output, lastErr
}

func selectNextAction[T any](plan ActionPlan[T], currentAction Action[T], direction string) (nextAction Action[T], err error) {
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
	// Wrap panic handling for safe running in a pipeline
	defer func() {
		if panicErr := recover(); panicErr != nil {
			logger.Errorf(ctx, "chain: panic occurred on running, caused by %v", panicErr)
			debug.PrintStack()

			output = input
			direction = Abort
			switch x := panicErr.(type) {
			case string:
				runError = errors.New(x)
			case error:
				runError = x
			default:
				runError = errors.New("unknown panic type")
			}
		}
	}()

	ctx = logger.WithRunnerDepth(ctx, action.Name())

	output, runError = action.Run(ctx, input)
	if runError != nil {
		return output, Error, runError
	}
	direction = Success
	if branchAction, isBranchAction := action.(BranchAction[T]); isBranchAction {
		direction, runError = branchAction.NextDirection(ctx, output)
	}

	return output, direction, runError
}
