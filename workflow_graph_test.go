package chain

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWorkflow_ValidateGraph(t *testing.T) {
	newAction := func(name string) Action[int] {
		return NewSimpleAction(
			name,
			func(_ context.Context, _ int) (int, error) { return 0, nil },
		)
	}

	type testCase struct {
		workflowConstructor func() *Workflow[int]
		isCyclic            bool
		isDisconnected      bool
	}

	testCases := map[string]testCase{
		"2 node cycle": {
			workflowConstructor: func() *Workflow[int] {
				action1 := newAction("action1")
				action2 := newAction("action2")

				workflow := NewWorkflow("workflow", action1, action2)
				// action1 and action2 makes cycle
				// (action1) -> action2
				// (action1) <- action2
				workflow.SetRunPlan(action1, SuccessOnlyPlan(action2))
				workflow.SetRunPlan(action2, SuccessOnlyPlan(action1))

				return workflow
			},
			isCyclic:       true,
			isDisconnected: false,
		},
		"3 node cycle": {
			workflowConstructor: func() *Workflow[int] {
				action1 := newAction("action1")
				action2 := newAction("action2")
				action3 := newAction("action3")

				workflow := NewWorkflow("workflow", action1, action2, action3)
				// (action1) -> action2 -> action3
				// (action1) <------------ action3
				workflow.SetRunPlan(action1, SuccessOnlyPlan(action2))
				workflow.SetRunPlan(action2, SuccessOnlyPlan(action3))
				workflow.SetRunPlan(action3, SuccessOnlyPlan(action1))

				return workflow
			},
			isCyclic:       true,
			isDisconnected: false,
		},
		"2 separate graph (disconnected)": {
			workflowConstructor: func() *Workflow[int] {
				action1 := newAction("action1")
				action2 := newAction("action2")
				action3 := newAction("action3")
				action4 := newAction("action4")

				workflow := NewWorkflow("workflow", action1, action2, action3, action4)
				// action1 -> action2 | action3 -> action4
				workflow.SetRunPlan(action1, SuccessOnlyPlan(action2))
				workflow.SetRunPlan(action2, TerminationPlan[int]())
				workflow.SetRunPlan(action3, SuccessOnlyPlan(action4))
				workflow.SetRunPlan(action4, TerminationPlan[int]())

				return workflow
			},
			isCyclic:       false,
			isDisconnected: true,
		},
		"valid dag but initAction is not entry node": {
			workflowConstructor: func() *Workflow[int] {
				action1 := newAction("action1")
				action2 := newAction("action2")
				action3 := newAction("action3")

				workflow := NewWorkflow("workflow", action2, action1, action3)
				// action1 -> (action2) -> action3
				workflow.SetRunPlan(action1, SuccessOnlyPlan(action2))
				workflow.SetRunPlan(action2, SuccessOnlyPlan(action3))
				workflow.SetRunPlan(action3, TerminationPlan[int]())

				return workflow
			},
			isCyclic:       false,
			isDisconnected: false,
		},
		"3 node cycle with branches": {
			workflowConstructor: func() *Workflow[int] {
				action1 := newAction("action1")
				action2 := newAction("action2")
				action3 := newAction("action3")

				workflow := NewWorkflow("workflow", action1, action2, action3)
				// (action1) ------------> action3
				// (action1) -> action2 -> action3
				// (action1) <- action2
				workflow.SetRunPlan(action1, DefaultPlan(action2, action3))
				workflow.SetRunPlan(action2, DefaultPlan(action3, action1))
				workflow.SetRunPlan(action3, TerminationPlan[int]())

				return workflow
			},
			isCyclic:       true,
			isDisconnected: false,
		},
		"valid dag but has 2 entry nodes": {
			workflowConstructor: func() *Workflow[int] {
				action1 := newAction("action1")
				action2 := newAction("action2")
				action3 := newAction("action3")
				action4 := newAction("action4")

				workflow := NewWorkflow("workflow", action1, action2, action3, action4)
				// (action1) -> action2 | action3 -> action4
				// (action1) ----------------------> action4
				workflow.SetRunPlan(action1, DefaultPlan(action2, action4))
				workflow.SetRunPlan(action2, TerminationPlan[int]())
				workflow.SetRunPlan(action3, SuccessOnlyPlan(action4))
				workflow.SetRunPlan(action4, TerminationPlan[int]())

				return workflow
			},
			isCyclic:       false,
			isDisconnected: false,
		},
		"valid dag but has 3 entry nodes": {
			workflowConstructor: func() *Workflow[int] {
				action1 := newAction("action1")
				action2 := newAction("action2")
				action3 := newAction("action3")
				action4 := newAction("action4")
				action5 := newAction("action5")

				workflow := NewWorkflow("workflow", action1, action2, action3, action4, action5)
				// (action1) -> action2 | action3 -> action4 <- action5
				// (action1) ----------------------> action4
				workflow.SetRunPlan(action1, DefaultPlan(action2, action4))
				workflow.SetRunPlan(action2, TerminationPlan[int]())
				workflow.SetRunPlan(action3, SuccessOnlyPlan(action4))
				workflow.SetRunPlan(action4, TerminationPlan[int]())
				workflow.SetRunPlan(action5, SuccessOnlyPlan(action4))

				return workflow
			},
			isCyclic:       false,
			isDisconnected: false,
		},
		"3 entry nodes, but one is disconnected": {
			workflowConstructor: func() *Workflow[int] {
				action1 := newAction("action1")
				action2 := newAction("action2")
				action3 := newAction("action3")
				action4 := newAction("action4")
				action5 := newAction("action5")

				workflow := NewWorkflow("workflow", action1, action2, action3, action4, action5)
				// (action1) -> action2 | action3 -> action4 | action5
				// (action1) ----------------------> action4
				workflow.SetRunPlan(action1, DefaultPlan(action2, action4))
				workflow.SetRunPlan(action2, TerminationPlan[int]())
				workflow.SetRunPlan(action3, SuccessOnlyPlan(action4))
				workflow.SetRunPlan(action4, TerminationPlan[int]())
				workflow.SetRunPlan(action5, TerminationPlan[int]())

				return workflow
			},
			isCyclic:       false,
			isDisconnected: true,
		},
		"valid branching with all directions": {
			workflowConstructor: func() *Workflow[int] {
				action1 := newAction("action1")
				action2 := newAction("action2")
				action3 := newAction("action3")

				workflow := NewWorkflow("workflow", action1, action2, action3)
				// (action1) -> action2
				// (action1) ------------> action3
				workflow.SetRunPlan(action1, DefaultPlan(action2, action3))
				workflow.SetRunPlan(action2, TerminationPlan[int]())
				workflow.SetRunPlan(action3, TerminationPlan[int]())

				return workflow
			},
			isCyclic:       false,
			isDisconnected: false,
		},
		"non cycle from initAction, but cycle in disconnected graph": {
			workflowConstructor: func() *Workflow[int] {
				action1 := newAction("action1")
				action2 := newAction("action2")
				action3 := newAction("action3")

				workflow := NewWorkflow("workflow", action1, action2, action3)
				// (action1) | action2 <-> action3
				workflow.SetRunPlan(action1, TerminationPlan[int]())
				workflow.SetRunPlan(action2, SuccessOnlyPlan(action3))
				workflow.SetRunPlan(action3, SuccessOnlyPlan(action2))

				return workflow
			},
			isCyclic:       true, // cyclic detected first
			isDisconnected: true,
		},
		"2 cycles": {
			workflowConstructor: func() *Workflow[int] {
				action1 := newAction("action1")
				action2 := newAction("action2")
				action3 := newAction("action3")
				action4 := newAction("action4")

				workflow := NewWorkflow("workflow", action1, action2, action3, action4)
				// (action1) -> action2 -> action3 -> action4
				// (action1) <-----------  action3 <- action4
				workflow.SetRunPlan(action1, SuccessOnlyPlan(action2))
				workflow.SetRunPlan(action2, SuccessOnlyPlan(action3))
				workflow.SetRunPlan(action3, DefaultPlan(action4, action1))
				workflow.SetRunPlan(action4, SuccessOnlyPlan(action3))
				workflow.SetRunPlan(action3, SuccessOnlyPlan(action4))

				return workflow
			},
			isCyclic:       true,
			isDisconnected: false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.NotPanics(t, func() {
				workflow := tc.workflowConstructor()

				err := workflow.ValidateGraph()
				switch {
				case tc.isCyclic:
					assert.Contains(t, err.Error(), "cycle")
				case tc.isDisconnected:
					assert.Contains(t, err.Error(), "disconnect")
				default:
					assert.NoError(t, err)
				}

				if err != nil {
					fmt.Println(err.Error())
				}
			})
		})
	}
}
