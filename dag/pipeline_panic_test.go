package dag

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPanicOnConfiguration(t *testing.T) {
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
				NewPipeline[string]("")
			},
			panics:  true,
			message: "pipeline must have a name",
		},
		"creating without action": {
			code: func() {
				NewPipeline[string]("Pipeline")
			},
			panics:  true,
			message: "no actions were described for creating pipeline",
		},
		"creating with terminate": {
			code: func() {
				action1 := &Blank{"action"}

				NewPipeline("Pipeline", action1, terminate)
			},
			panics:  true,
			message: "do not set terminate as a member",
		},
		"creating with duplicate actions": {
			code: func() {
				action1 := &Blank{"action"}

				NewPipeline("Pipeline", action1, action1)
			},
			panics:  true,
			message: "duplicate action specified on actions argument 2",
		},
		"set action plan for terminate": {
			code: func() {
				action1 := &Blank{"action"}
				pipeline := NewPipeline("Pipeline", action1)

				pipeline.SetRunPlan(terminate, terminationPlan)
			},
			panics:  true,
			message: "cannot set plan for terminate",
		},
		"set action plan for non-member": {
			code: func() {
				action1, nonMember := &Blank{"member"}, &Blank{"non-member"}
				pipeline := NewPipeline("Pipeline", action1)

				pipeline.SetRunPlan(nonMember, ActionPlan[string]{
					Success: action1,
				})
			},
			panics:  true,
			message: "`non-member` is not a member of this pipeline",
		},
		"set action plan directing non-member": {
			code: func() {
				action1, nonMember := &Blank{"member"}, &Blank{"non-member"}
				pipeline := NewPipeline("Pipeline", action1)

				pipeline.SetRunPlan(action1, ActionPlan[string]{
					Success: nonMember,
				})
			},
			panics:  true,
			message: "setting plan from `member` directing `success` to non-member `non-member`",
		},
		"set action plan with self loop": {
			code: func() {
				action1 := &Blank{"member"}
				pipeline := NewPipeline("Pipeline", action1)

				pipeline.SetRunPlan(action1, ActionPlan[string]{
					Success: action1,
				})
			},
			panics:  true,
			message: "setting self loop plan with `member` directing `success`",
		},
		"not configuring manual plans": {
			code: func() {
				action1 := &Blank{"action"}
				pipeline := NewPipeline("Pipeline", action1)

				_ = pipeline
			},
			panics: false,
		},
		"skipping some default directions": {
			code: func() {
				action1 := &Blank{"action"}
				pipeline := NewPipeline("Pipeline", action1)

				pipeline.SetRunPlan(action1, ActionPlan[string]{
					Success: terminate,
					// Not configuring Error, Abort
				})
			},
			panics: false,
		},
		"skipping plan description": {
			code: func() {
				action1 := &Blank{"action"}
				pipeline := NewPipeline("Pipeline", action1)

				pipeline.SetRunPlan(action1, terminationPlan) // terminationPlan is nil
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

type Blank struct{ name string }

func (b Blank) Name() string { return b.name }
func (Blank) Run(_ context.Context, _ string) (string, string, error) {
	return "", Abort, nil
}
