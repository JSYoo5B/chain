package v1

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
)

type Pipeline[T any] struct {
	name       string
	runPlans   map[Action[T]]ActionPlan[T]
	initAction Action[T]
}

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
		p.runPlans[action] = SuccessOnlyPlan(nextAction)
	}

	return p
}

func (p *Pipeline[T]) SetRunPlan(currentAction Action[T], plan ActionPlan[T]) {
	if currentAction == nil {
		panic(errors.New("cannot set plan for terminate"))
	}
	if _, exists := p.runPlans[currentAction]; !exists {
		panic(fmt.Errorf("`%s` is not a member of this pipeline", currentAction.Name()))
	}

	// When given plan is nil, make currentAction to terminate on any cases
	if plan == nil {
		plan = ActionPlan[T]{}
	}

	// Set next action to terminate when allowed directions were not specified in plan
	terminate := Terminate[T]()
	availableDirections := append(currentAction.Directions(), Success, Error, Abort)
	for _, direction := range availableDirections {
		if _, exists := plan[direction]; !exists {
			plan[direction] = terminate
		}
	}

	// Validate given plan with members
	var panicMsg error
	for direction, nextAction := range plan {
		if nextAction == terminate {
			continue
		}

		// If the direction is not in currentAction's valid directions, panic
		if !contains(currentAction.Directions(), direction) {
			panicMsg = fmt.Errorf("`%s` does not support direction `%s`", currentAction.Name(), direction)
		}

		// Check the nextAction is valid. (check non-member, or self loop)
		if !p.isMember(nextAction) {
			panicMsg = fmt.Errorf("setting plan from `%s` directing `%s` to non-member `%s`", currentAction.Name(), direction, nextAction.Name())
		} else if nextAction == currentAction {
			panicMsg = fmt.Errorf("setting self loop plan with `%s` directing `%s`", currentAction.Name(), direction)
		}

		if panicMsg != nil {
			panic(panicMsg)
		}
	}

	p.runPlans[currentAction] = plan
}

func (p *Pipeline[T]) Name() string         { return p.name }
func (p *Pipeline[T]) Directions() []string { return []string{Success, Error, Abort} }
func (p *Pipeline[T]) Run(ctx context.Context, input T) (output T, direction string, err error) {
	if len(p.runPlans) == 1 {
		return runAction(p.initAction, ctx, input)
	}

	return p.RunAt(p.initAction, ctx, input)
}

const parentRunner = "PipelineParentRunner"

func (p *Pipeline[T]) RunAt(initAction Action[T], ctx context.Context, input T) (output T, direction string, lastErr error) {
	if _, exists := p.runPlans[initAction]; !exists {
		return input, Error, errors.New("given initAction is not registered on constructor")
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
		runErr        error
		selectErr     error
	)
	logrus.Debugf("%s: Start running with `%s`", runnerName, initAction.Name())
	for currentAction = initAction; currentAction != nil; currentAction = nextAction {
		output, direction, runErr = runAction(currentAction, ctx, input)

		nextAction, selectErr = p.selectNextAction(currentAction, direction)
		if selectErr != nil {
			logrus.Error(selectErr)
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

	return output, direction, lastErr
}

func (p *Pipeline[T]) selectNextAction(currentAction Action[T], direction string) (nextAction Action[T], err error) {
	terminate := Terminate[T]()
	plan, exist := p.runPlans[currentAction]
	if !exist || plan == nil {
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

func (p *Pipeline[T]) isMember(action Action[T]) bool {
	_, exists := p.runPlans[action]
	return exists
}
