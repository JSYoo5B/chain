package chain

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWorkflow_Configure(t *testing.T) {
	newAction := func(name string) Action[string] {
		return NewSimpleAction(
			name,
			func(_ context.Context, _ string) (string, error) { return "", nil },
		)
	}

	terminate := Terminate[string]()
	terminationPlan := TerminationPlan[string]()
	type testCase struct {
		code    func()
		panics  bool
		message string
	}

	testCases := map[string]testCase{
		"creating without name": {
			code: func() {
				NewWorkflow[string]("")
			},
			panics:  true,
			message: "workflow must have a name",
		},
		"creating without action": {
			code: func() {
				NewWorkflow[string]("Workflow")
			},
			panics:  true,
			message: "no actions were described for creating workflow",
		},
		"creating with terminate": {
			code: func() {
				action1 := newAction("action")

				NewWorkflow("Workflow", action1, terminate)
			},
			panics:  true,
			message: "do not set terminate as a member",
		},
		"creating with duplicate actions": {
			code: func() {
				action1 := newAction("action")

				NewWorkflow("Workflow", action1, action1)
			},
			panics:  true,
			message: "duplicate action specified on actions argument 2",
		},
		"set action plan for terminate": {
			code: func() {
				action1 := newAction("action")
				workflow := NewWorkflow("Workflow", action1)

				workflow.SetRunPlan(terminate, terminationPlan)
			},
			panics:  true,
			message: "cannot set plan for terminate",
		},
		"set action plan for unsupported direction": {
			code: func() {
				action1, action2 := newAction("action1"), newAction("action2")
				workflow := NewWorkflow("Workflow", action1, action2)

				workflow.SetRunPlan(action1, ActionPlan[string]{
					"unsupported": action2,
				})
			},
			panics:  true,
			message: "`action1` does not support direction `unsupported`",
		},
		"set action plan for non-member": {
			code: func() {
				action1, nonMember := newAction("member"), newAction("non-member")
				workflow := NewWorkflow("Workflow", action1)

				workflow.SetRunPlan(nonMember, SuccessOnlyPlan(action1))
			},
			panics:  true,
			message: "`non-member` is not a member of this workflow",
		},
		"set action plan directing non-member": {
			code: func() {
				action1, nonMember := newAction("member"), newAction("non-member")
				workflow := NewWorkflow("Workflow", action1)

				workflow.SetRunPlan(action1, SuccessOnlyPlan(nonMember))
			},
			panics:  true,
			message: "setting plan from `member` directing `success` to non-member `non-member`",
		},
		"set action plan with self loop": {
			code: func() {
				action1 := newAction("member")
				workflow := NewWorkflow("Workflow", action1)

				workflow.SetRunPlan(action1, SuccessOnlyPlan(action1))
			},
			panics:  true,
			message: "setting self loop plan with `member` directing `success`",
		},
		"not configuring manual plans": {
			code: func() {
				action1 := newAction("action")
				workflow := NewWorkflow("Workflow", action1)

				_ = workflow
			},
			panics: false,
		},
		"skipping some default directions": {
			code: func() {
				action1 := newAction("action")
				workflow := NewWorkflow("Workflow", action1)

				workflow.SetRunPlan(action1, ActionPlan[string]{
					Success: terminate,
					// Not configuring Failure, Abort
				})
			},
			panics: false,
		},
		"skipping plan description": {
			code: func() {
				action1 := newAction("action")
				workflow := NewWorkflow("Workflow", action1)

				workflow.SetRunPlan(action1, terminationPlan) // terminationPlan is nil
			},
			panics: false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			if tc.panics {
				assert.PanicsWithError(t, tc.message, tc.code)
			} else {
				assert.NotPanics(t, tc.code)
			}
		})
	}
}
