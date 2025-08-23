package chain

import (
	"errors"
	"fmt"
)

// Workflow represents a sequence of Actions that are executed in a structured flow.
// It executes each of its constituent Actions in sequence, with each Action following its
// own Run method. The flow proceeds based on the defined structure of the Workflow,
// allowing flexible and organized execution of actions to build workflows that can be
// as simple or complex as needed.
//
// Workflow implements the Action interface, meaning it can be treated as an Action itself.
// This allows Workflows to be composed hierarchically, enabling more complex workflows by nesting
// Workflows within other Workflows.
type Workflow[T any] struct {
	name       string
	runPlans   map[Action[T]]ActionPlan[T]
	initAction Action[T]
}

// NewWorkflow creates a new Workflow by taking a series of Actions as its members.
// These Actions will be executed sequentially in the order they are provided, with the output
// of one Action being passed as input to the next, forming a unidirectional flow of execution.
func NewWorkflow[T any](name string, memberActions ...Action[T]) *Workflow[T] {
	if name == "" {
		panic(errors.New("workflow must have a name"))
	}
	if len(memberActions) == 0 {
		panic(errors.New("no actions were described for creating workflow"))
	}

	p := &Workflow[T]{
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

// SetRunPlan updates the execution flow for the given currentAction in the Workflow
// by associating it with a specified ActionPlan. The currentAction will be validated
// to ensure it is a member of the Workflow. The ActionPlan defines the directions
// (such as Success, Error, Abort) and their corresponding next actions in the execution flow.
//
// If the currentAction is nil or not part of the Workflow, a panic will occur.
// The plan can be nil, in which case the currentAction will be set to terminate
// for any direction not explicitly specified in the plan. If a direction is
// encountered in the plan that is not valid for the currentAction, or if it
// leads to an invalid action, another panic will occur.
//
// Additionally, self-loops are not allowed in the plan. If the next action for
// a direction is the current action itself, a panic will be triggered.
func (w *Workflow[T]) SetRunPlan(currentAction Action[T], plan ActionPlan[T]) {
	if currentAction == nil {
		panic(errors.New("cannot set plan for terminate"))
	} else if !isMemberActionInWorkflow(currentAction, w) {
		panic(fmt.Errorf("`%s` is not a member of this workflow", currentAction.Name()))
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
	for _, direction := range availableDirections {
		if _, exists := plan[direction]; !exists {
			plan[direction] = terminate
		}
	}

	// Validate a given plan with members
	var err error
	for direction, nextAction := range plan {
		if nextAction == terminate {
			continue
		}

		// If the direction is not in the currentAction's valid directions, panic
		if !contains(availableDirections, direction) {
			err = fmt.Errorf("`%s` does not support direction `%s`", currentAction.Name(), direction)
		} else if !isMemberActionInWorkflow(nextAction, w) {
			err = fmt.Errorf("setting plan from `%s` directing `%s` to non-member `%s`", currentAction.Name(), direction, nextAction.Name())
		} else if nextAction == currentAction {
			err = fmt.Errorf("setting self loop plan with `%s` directing `%s`", currentAction.Name(), direction)
		}

		if err != nil {
			panic(err)
		}
	}

	w.runPlans[currentAction] = plan
}

// Name provides the identifier of this Workflow.
func (w *Workflow[T]) Name() string { return w.name }

func contains(directions []string, direction string) bool {
	for _, dir := range directions {
		if dir == direction {
			return true
		}
	}
	return false
}

func isMemberActionInWorkflow[T any](action Action[T], p *Workflow[T]) bool {
	_, exists := p.runPlans[action]
	return exists
}
