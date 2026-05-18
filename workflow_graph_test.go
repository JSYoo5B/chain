package chain

import (
	"context"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorkflow_ValidateGraph(t *testing.T) {
	newActions := func(count int) []Action[int] {
		actions := make([]Action[int], count)
		for i := range count {
			name := strconv.Itoa(i)
			actions[i] = NewSimpleAction(
				name,
				func(_ context.Context, _ int) (int, error) { return 0, nil },
			)
		}
		return actions
	}

	type testCase struct {
		workflowConstructor func() *Workflow[int]
		expectedErrContains []string
		repetitions         int
	}

	testCases := map[string]testCase{
		"valid: single node graph": {
			workflowConstructor: func() *Workflow[int] {
				actions := newActions(1)
				workflow := NewWorkflow("workflow", actions...)

				// (0)
				workflow.SetRunPlan(actions[0], TerminationPlan[int]())

				return workflow
			},
		},
		"valid: linear graph when initAction is not entry node": {
			workflowConstructor: func() *Workflow[int] {
				actions := newActions(3)
				workflow := NewWorkflow("workflow", actions[1], actions[0], actions[2])

				// 0 -> (1) -> 2
				workflow.SetRunPlan(actions[0], SuccessOnlyPlan(actions[1]))
				workflow.SetRunPlan(actions[1], SuccessOnlyPlan(actions[2]))
				workflow.SetRunPlan(actions[2], TerminationPlan[int]())

				return workflow
			},
		},
		"valid: branch graph": {
			workflowConstructor: func() *Workflow[int] {
				actions := newActions(3)
				workflow := NewWorkflow("workflow", actions...)

				// (0) -> 1
				// (0) ----> 2
				workflow.SetRunPlan(actions[0], DefaultPlan(actions[1], actions[2]))
				workflow.SetRunPlan(actions[1], TerminationPlan[int]())
				workflow.SetRunPlan(actions[2], TerminationPlan[int]())

				return workflow
			},
		},
		"valid: multiple directions point to same node": {
			workflowConstructor: func() *Workflow[int] {
				actions := newActions(2)
				workflow := NewWorkflow("workflow", actions...)

				// (0) -success-> 1
				// (0) -failure-> 1
				workflow.SetRunPlan(actions[0], DefaultPlan(actions[1], actions[1]))
				workflow.SetRunPlan(actions[1], TerminationPlan[int]())

				return workflow
			},
		},
		"valid: branch and merge graph": {
			workflowConstructor: func() *Workflow[int] {
				actions := newActions(4)
				workflow := NewWorkflow("workflow", actions...)

				// (0) -> 1 ----> 3
				// (0) ----> 2 -> 3
				workflow.SetRunPlan(actions[0], DefaultPlan(actions[1], actions[2]))
				workflow.SetRunPlan(actions[1], SuccessOnlyPlan(actions[3]))
				workflow.SetRunPlan(actions[2], SuccessOnlyPlan(actions[3]))
				workflow.SetRunPlan(actions[3], TerminationPlan[int]())

				return workflow
			},
		},
		"valid: weakly connected graph by common upstream action": {
			workflowConstructor: func() *Workflow[int] {
				actions := newActions(4)
				workflow := NewWorkflow("workflow", actions...)

				// 1 -> (0)
				// 1 ----> 2
				// 1 --------> 3
				workflow.SetRunPlan(actions[0], TerminationPlan[int]())
				workflow.SetRunPlan(actions[1], RunPlan[int]{
					Success: actions[0],
					Failure: actions[2],
					Abort:   actions[3],
				})
				workflow.SetRunPlan(actions[2], TerminationPlan[int]())
				workflow.SetRunPlan(actions[3], TerminationPlan[int]())

				return workflow
			},
			repetitions: 100,
		},
		"invalid: cycle": {
			workflowConstructor: func() *Workflow[int] {
				actions := newActions(3)
				workflow := NewWorkflow("workflow", actions...)

				// (0) -> 1 -> 2
				// (0) <------ 2
				workflow.SetRunPlan(actions[0], SuccessOnlyPlan(actions[1]))
				workflow.SetRunPlan(actions[1], SuccessOnlyPlan(actions[2]))
				workflow.SetRunPlan(actions[2], SuccessOnlyPlan(actions[0]))

				return workflow
			},
			expectedErrContains: []string{"cycle"},
		},
		"invalid: cycle through branch": {
			workflowConstructor: func() *Workflow[int] {
				actions := newActions(3)
				workflow := NewWorkflow("workflow", actions...)

				// (0) -> 1 -> 2
				// (0) <-- 1
				workflow.SetRunPlan(actions[0], DefaultPlan(actions[1], actions[2]))
				workflow.SetRunPlan(actions[1], DefaultPlan(actions[2], actions[0]))
				workflow.SetRunPlan(actions[2], TerminationPlan[int]())

				return workflow
			},
			expectedErrContains: []string{"cycle"},
		},
		"invalid: cycle outside initAction directed path": {
			workflowConstructor: func() *Workflow[int] {
				actions := newActions(3)
				workflow := NewWorkflow("workflow", actions...)

				// 1 -> (0)
				// 1 ----> 2
				// 1 <---- 2
				workflow.SetRunPlan(actions[0], TerminationPlan[int]())
				workflow.SetRunPlan(actions[1], DefaultPlan(actions[0], actions[2]))
				workflow.SetRunPlan(actions[2], SuccessOnlyPlan(actions[1]))

				return workflow
			},
			expectedErrContains: []string{"cycle"},
		},
		"invalid: disconnected graph": {
			workflowConstructor: func() *Workflow[int] {
				actions := newActions(4)
				workflow := NewWorkflow("workflow", actions...)

				// 0 -> 1 | 2 -> 3
				workflow.SetRunPlan(actions[0], SuccessOnlyPlan(actions[1]))
				workflow.SetRunPlan(actions[1], TerminationPlan[int]())
				workflow.SetRunPlan(actions[2], SuccessOnlyPlan(actions[3]))
				workflow.SetRunPlan(actions[3], TerminationPlan[int]())

				return workflow
			},
			expectedErrContains: []string{"disconnected graph"},
		},
		"invalid: disconnected graph with isolated node": {
			workflowConstructor: func() *Workflow[int] {
				actions := newActions(3)
				workflow := NewWorkflow("workflow", actions...)

				// (0) -> 1 | 2
				workflow.SetRunPlan(actions[0], SuccessOnlyPlan(actions[1]))
				workflow.SetRunPlan(actions[1], TerminationPlan[int]())
				workflow.SetRunPlan(actions[2], TerminationPlan[int]())

				return workflow
			},
			expectedErrContains: []string{"disconnected graph"},
		},
		"invalid: disconnected graph with cycle": {
			workflowConstructor: func() *Workflow[int] {
				actions := newActions(3)
				workflow := NewWorkflow("workflow", actions...)

				// (0) | 1 -> 2
				//     | 1 <- 2
				workflow.SetRunPlan(actions[0], TerminationPlan[int]())
				workflow.SetRunPlan(actions[1], SuccessOnlyPlan(actions[2]))
				workflow.SetRunPlan(actions[2], SuccessOnlyPlan(actions[1]))

				return workflow
			},
			expectedErrContains: []string{"cycle", "disconnected graph"},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			repetitions := tc.repetitions
			if repetitions == 0 {
				repetitions = 1
			}

			for i := 0; i < repetitions; i++ {
				workflow := tc.workflowConstructor()

				err := workflow.ValidateGraph()
				if len(tc.expectedErrContains) == 0 {
					if !assert.NoError(t, err) {
						break
					}
					continue
				}

				if assert.Error(t, err) {
					assertErrorContainsAny(t, err, tc.expectedErrContains)
				}
			}
		})
	}
}

func assertErrorContainsAny(t *testing.T, err error, expectedSubstrings []string) {
	t.Helper()

	for _, expectedSubstring := range expectedSubstrings {
		if strings.Contains(err.Error(), expectedSubstring) {
			return
		}
	}

	assert.Failf(t, "unexpected error", "expected error %q to contain one of %v", err.Error(), expectedSubstrings)
}
