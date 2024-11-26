package railway

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"runtime/debug"
)

type Pipeline[T any] struct {
	name       string
	runPlans   map[Action[T]]ActionPlan[T]
	initAction Action[T]
}

// NewPipeline initializes a new Pipeline with the provided name and actions.
// It sets up the default ActionPlan for each Action in the pipeline, linking
// each action to the next in sequence.
// If no next action is specified, it links to the terminate action by default.
// It ensures that no duplicate actions or the terminate action itself are
// included in the pipeline, throwing a panic if either is detected.
// This method establishes the sequential flow of actions by setting
// the next action for each action in the plan.
func NewPipeline[T any](name string, memberActions ...Action[T]) *Pipeline[T] {
	if name == "" {
		panic(errors.New("pipeline must have a name"))
	}
	if len(memberActions) == 0 {
		panic(errors.New("no actions were described for creating pipeline"))
	}

	p := &Pipeline[T]{
		name:       name,
		runPlans:   map[Action[T]]ActionPlan[T]{},
		initAction: memberActions[0],
	}

	terminate := Terminate[T]()
	for i, action := range memberActions {
		if action == terminate {
			panic(errors.New("do not set terminate as a member"))
		}
		if _, exists := p.runPlans[action]; exists {
			panic(fmt.Errorf("duplicate action specified on actions argument %d", i+1))
		}

		nextAction := terminate
		if i+1 < len(memberActions) {
			nextAction = memberActions[i+1]
		}

		defaultPlan := ActionPlan[T]{}
		availableDirections := []string{Success, Error, Abort}
		if branchAction, isBranchAction := action.(BranchAction[T]); isBranchAction {
			availableDirections = append(availableDirections, branchAction.Directions()...)
		}
		for _, direction := range availableDirections {
			if _, exists := defaultPlan[direction]; !exists {
				defaultPlan[direction] = terminate
			}
		}
		defaultPlan[Success] = nextAction
		p.runPlans[action] = defaultPlan
	}

	return p
}

// SetRunPlan sets the action plan for a given action in the pipeline, linking each direction
// (Success, Error, Abort and other custom branching directions) to a subsequent action.
// It ensures that only valid directions supported by the currentAction are used.
// Any direction that the current action does not support leads to a panic.
// Additionally, if the next action is not part of the pipeline, or if there's
// a self-loop, a panic will occur.
// If no direction is planned, it defaults to terminating with the `terminate` action.
func (p *Pipeline[T]) SetRunPlan(currentAction Action[T], plan ActionPlan[T]) {
	if currentAction == nil {
		panic(errors.New("cannot set plan for terminate"))
	} else if !isMemberActionInPipeline(currentAction, p) {
		panic(fmt.Errorf("`%s` is not a member of this pipeline", currentAction.Name()))
	}

	// When given plan is nil, make currentAction to terminate on any cases
	if plan == nil {
		plan = ActionPlan[T]{}
	}

	// Set next action to terminate when allowed directions were not specified in plan
	terminate := Terminate[T]()
	availableDirections := []string{Success, Error, Abort}
	if branchAction, isBranchAction := currentAction.(BranchAction[T]); isBranchAction {
		availableDirections = append(availableDirections, branchAction.Directions()...)
	}
	for _, direction := range append(availableDirections, Success, Error, Abort) {
		if _, exists := plan[direction]; !exists {
			plan[direction] = terminate
		}
	}

	// Validate given plan with members
	var err error
	for direction, nextAction := range plan {
		if nextAction == terminate {
			continue
		}

		// If the direction is not in currentAction's valid directions, panic
		if !contains(availableDirections, direction) {
			err = fmt.Errorf("`%s` does not support direction `%s`", currentAction.Name(), direction)
		} else if !isMemberActionInPipeline(nextAction, p) {
			err = fmt.Errorf("setting plan from `%s` directing `%s` to non-member `%s`", currentAction.Name(), direction, nextAction.Name())
		} else if nextAction == currentAction {
			err = fmt.Errorf("setting self loop plan with `%s` directing `%s`", currentAction.Name(), direction)
		}

		if err != nil {
			panic(err)
		}
	}

	p.runPlans[currentAction] = plan
}

// Name returns the name of the Pipeline, which is a distinguishable identifier for the pipeline.
func (p *Pipeline[T]) Name() string { return p.name }

// Run executes its member Action within the Pipeline sequentially
// from the initAction to termination, following the specified ActionPlan.
// The initAction refers to the first Action
// in the memberActions provided as an argument to NewPipeline.
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

	runnerName := p.name
	if parentName := ctx.Value(parentRunner); parentName != nil {
		runnerName = parentName.(string) + "/" + p.name
	}
	ctx = context.WithValue(ctx, parentRunner, runnerName)

	var (
		terminate     = Terminate[T]()
		currentAction Action[T]
		nextAction    Action[T]
		direction     string
		runErr        error
		selectErr     error
	)
	logrus.Debugf("%s: Start running with `%s`", runnerName, initAction.Name())
	for currentAction = initAction; currentAction != nil; currentAction = nextAction {
		output, direction, runErr = runAction(currentAction, ctx, input)

		nextAction, selectErr = selectNextAction(p.runPlans[currentAction], currentAction, direction)
		if selectErr != nil {
			logrus.Error(selectErr)
			direction = Abort
			lastErr = selectErr
			break
		}

		nextActionName := "termination"
		if nextAction != terminate {
			nextActionName = nextAction.Name()
		}
		logrus.Debugf("%s: `%s` directs `%s`, selecting `%s`", runnerName, currentAction.Name(), direction, nextActionName)

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

const parentRunner = "PipelineParentRunner"

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

func contains(directions []string, direction string) bool {
	for _, dir := range directions {
		if dir == direction {
			return true
		}
	}
	return false
}

func isMemberActionInPipeline[T any](action Action[T], p *Pipeline[T]) bool {
	_, exists := p.runPlans[action]
	return exists
}

func runAction[T any](action Action[T], ctx context.Context, input T) (output T, direction string, runError error) {
	// Wrap panic handling for safe running in pipeline
	defer func() {
		if panicErr := recover(); panicErr != nil {
			logrus.Errorf("%s: panic occurred on running, caused by %s", action.Name(), panicErr)
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
