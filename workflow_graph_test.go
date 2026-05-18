package chain

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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
		expectedErrContains []string
		repetitions         int
	}

	testCases := map[string]testCase{
		"valid linear graph when initAction is not entry node": {
			workflowConstructor: func() *Workflow[int] {
				a1 := newAction("1")
				a2 := newAction("2")
				a3 := newAction("3")

				workflow := NewWorkflow("workflow", a2, a1, a3)
				// 1 -> (2) -> 3
				workflow.SetRunPlan(a1, SuccessOnlyPlan(a2))
				workflow.SetRunPlan(a2, SuccessOnlyPlan(a3))
				workflow.SetRunPlan(a3, TerminationPlan[int]())

				return workflow
			},
		},
		"valid branch graph": {
			workflowConstructor: func() *Workflow[int] {
				a1 := newAction("1")
				a2 := newAction("2")
				a3 := newAction("3")

				workflow := NewWorkflow("workflow", a1, a2, a3)
				// (1) -> 2
				// (1) ----> 3
				workflow.SetRunPlan(a1, DefaultPlan(a2, a3))
				workflow.SetRunPlan(a2, TerminationPlan[int]())
				workflow.SetRunPlan(a3, TerminationPlan[int]())

				return workflow
			},
		},
		"cycle": {
			workflowConstructor: func() *Workflow[int] {
				a1 := newAction("1")
				a2 := newAction("2")
				a3 := newAction("3")

				workflow := NewWorkflow("workflow", a1, a2, a3)
				// (1) -> 2 -> 3
				// (1) <------ 3
				workflow.SetRunPlan(a1, SuccessOnlyPlan(a2))
				workflow.SetRunPlan(a2, SuccessOnlyPlan(a3))
				workflow.SetRunPlan(a3, SuccessOnlyPlan(a1))

				return workflow
			},
			expectedErrContains: []string{"cycle"},
		},
		"cycle through branch": {
			workflowConstructor: func() *Workflow[int] {
				a1 := newAction("1")
				a2 := newAction("2")
				a3 := newAction("3")

				workflow := NewWorkflow("workflow", a1, a2, a3)
				// (1) -> 2 -> 3
				// (1) <-- 2
				workflow.SetRunPlan(a1, DefaultPlan(a2, a3))
				workflow.SetRunPlan(a2, DefaultPlan(a3, a1))
				workflow.SetRunPlan(a3, TerminationPlan[int]())

				return workflow
			},
			expectedErrContains: []string{"cycle"},
		},
		"disconnected graph": {
			workflowConstructor: func() *Workflow[int] {
				a1 := newAction("1")
				a2 := newAction("2")
				a3 := newAction("3")
				a4 := newAction("4")

				workflow := NewWorkflow("workflow", a1, a2, a3, a4)
				// 1 -> 2 | 3 -> 4
				workflow.SetRunPlan(a1, SuccessOnlyPlan(a2))
				workflow.SetRunPlan(a2, TerminationPlan[int]())
				workflow.SetRunPlan(a3, SuccessOnlyPlan(a4))
				workflow.SetRunPlan(a4, TerminationPlan[int]())

				return workflow
			},
			expectedErrContains: []string{"disconnect"},
		},
		"disconnected graph with cycle": {
			workflowConstructor: func() *Workflow[int] {
				a1 := newAction("1")
				a2 := newAction("2")
				a3 := newAction("3")

				workflow := NewWorkflow("workflow", a1, a2, a3)
				// (1) | 2 -> 3
				//     | 2 <- 3
				workflow.SetRunPlan(a1, TerminationPlan[int]())
				workflow.SetRunPlan(a2, SuccessOnlyPlan(a3))
				workflow.SetRunPlan(a3, SuccessOnlyPlan(a2))

				return workflow
			},
			expectedErrContains: []string{"cycle", "disconnect"},
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
