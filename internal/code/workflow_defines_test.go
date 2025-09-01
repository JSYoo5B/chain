package code

import (
	"github.com/JSYoo5B/chain"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestWorkflowEdge_ActionPlanCode(t *testing.T) {
	type testCase struct {
		edge     WorkflowEdge
		expected string
	}

	testCases := map[string]testCase{
		"no plans": {
			edge: WorkflowEdge{
				WorkType: "int",
				Plan:     nil,
			},
			expected: "chain.TerminationPlan[int]()",
		},
		"success only plan": {
			edge: WorkflowEdge{
				WorkType: "int",
				Plan:     map[string]string{chain.Success: "odd2"},
			},
			expected: "chain.SuccessOnlyPlan(odd2)",
		},
		"default plan": {
			edge: WorkflowEdge{
				WorkType: "int",
				Plan: map[string]string{
					chain.Success: "odd2",
					chain.Error:   "error",
				},
			},
			expected: "chain.DefaultPlan(odd2,error)",
		},
		"default plan with abort": {
			edge: WorkflowEdge{
				WorkType: "int",
				Plan: map[string]string{
					chain.Success: "odd2",
					chain.Error:   "error",
					chain.Abort:   "abort",
				},
			},
			expected: "chain.DefaultPlanWithAbort(odd2,error,abort)",
		},
		"custom branch plan": {
			edge: WorkflowEdge{
				WorkType: "int",
				Plan: map[string]string{
					"even": "even",
					"odd":  "odd1",
				},
			},
			expected: `chain.ActionPlan[int]{"even":even,"odd":odd1,}`,
		},
	}

	removeAllBlank := func(s string) string {
		re := regexp.MustCompile(`\s+`)
		return re.ReplaceAllString(s, "")
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			actual := tc.edge.ActionPlanCode()

			assert.Equal(t, tc.expected, removeAllBlank(actual))
		})
	}
}
