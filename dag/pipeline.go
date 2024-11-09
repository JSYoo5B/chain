package dag

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
)

type Pipeline[T any] struct {
	name       string
	initAction Action[T]
	planMap    map[Action[T]]ActionPlan[T]
}

func NewPipeline[T any](name string, actions ...Action[T]) *Pipeline[T] {
	p := &Pipeline[T]{
		name:    name,
		planMap: map[Action[T]]ActionPlan[T]{},
	}

	if len(actions) == 0 {
		panic("no actions were described for creating pipeline")
	}

	p.initAction = actions[0]
	terminate := TerminateAction[T]()
	for i, action := range actions {
		nextAction := terminate
		if i+1 < len(actions) {
			nextAction = actions[i+1]
		}

		defaultPlan := ActionPlan[T]{
			Success: nextAction,
			Error:   terminate,
			Abort:   terminate,
		}

		if _, exists := p.planMap[action]; exists {
			panic(fmt.Sprintf("duplicate action specified on actions argument %d", i+1))
		}

		p.planMap[action] = defaultPlan
	}

	return p
}

func (p *Pipeline[T]) Name() string { return p.name }

func (p *Pipeline[T]) SetActionPlan(currentAction Action[T], plan ActionPlan[T]) {
	if _, exists := p.planMap[currentAction]; !exists {
		panic("given action is not registered on constructor")
	}

	// When given plan is nil, make currentAction to terminate on any cases
	if plan == nil {
		plan = ActionPlan[T]{}
	}

	// Set next action to terminate when default directions were not planned
	terminate := TerminateAction[T]()
	if _, exists := plan[Success]; !exists {
		plan[Success] = terminate
	}
	if _, exists := plan[Error]; !exists {
		plan[Error] = terminate
	}
	if _, exists := plan[Abort]; !exists {
		plan[Abort] = terminate
	}

	p.planMap[currentAction] = plan
}

func (p *Pipeline[T]) Run(ctx context.Context, input T) (output T, direction string, err error) {
	return p.RunAt(p.initAction, ctx, input)
}

const parentRunner = "PipelineParentRunner"

func (p *Pipeline[T]) RunAt(initAction Action[T], ctx context.Context, input T) (output T, direction string, runError error) {
	if _, exists := p.planMap[initAction]; !exists {
		return input, Error, errors.New("given initAction is not registered on constructor")
	}

	runnerName := p.name
	if parentName := ctx.Value(parentRunner); parentName != nil {
		runnerName = parentName.(string) + "/" + p.name
	}
	ctx = context.WithValue(ctx, parentRunner, runnerName)

	var (
		terminate     = TerminateAction[T]()
		currentAction Action[T]
		nextAction    Action[T]
		selectErr     error
	)
	logrus.Debugf("%s: Start running with %s", runnerName, initAction.Name())
	for currentAction = initAction; currentAction != nil; currentAction = nextAction {
		output, direction, runError = runAction(currentAction, ctx, input)

		nextAction, selectErr = p.selectNextAction(currentAction, direction)
		if selectErr != nil {
			logrus.Errorf("failed to select next Action: %s", selectErr.Error())
			break
		}

		nextActionName := "termination"
		if nextAction != terminate {
			nextActionName = nextAction.Name()
		}
		logrus.Debugf("%s: %s directs '%s', selecting %s", runnerName, currentAction.Name(), direction, nextActionName)

		input = output
	}

	return output, direction, runError
}

func (p *Pipeline[T]) selectNextAction(currentAction Action[T], direction string) (nextAction Action[T], err error) {
	terminate := TerminateAction[T]()
	plan, exist := p.planMap[currentAction]
	if !exist || plan == nil {
		return terminate, fmt.Errorf("no action plan found for Action[%s]", currentAction.Name())
	}
	if nextAction, exist = plan[direction]; !exist {
		return terminate, fmt.Errorf("no action plan from Action[%s] directing %s", currentAction.Name(), direction)
	}

	return nextAction, nil
}
