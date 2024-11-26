package chain

import (
	"context"
	"github.com/sirupsen/logrus"
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
		"set action plan for unsupported direction": {
			code: func() {
				action1, action2 := &Blank{"action1"}, &Blank{"action2"}
				pipeline := NewPipeline("Pipeline", action1, action2)

				pipeline.SetRunPlan(action1, ActionPlan[string]{
					"unsupported": action2,
				})
			},
			panics:  true,
			message: "`action1` does not support direction `unsupported`",
		},
		"set action plan for non-member": {
			code: func() {
				action1, nonMember := &Blank{"member"}, &Blank{"non-member"}
				pipeline := NewPipeline("Pipeline", action1)

				pipeline.SetRunPlan(nonMember, SuccessOnlyPlan(action1))
			},
			panics:  true,
			message: "`non-member` is not a member of this pipeline",
		},
		"set action plan directing non-member": {
			code: func() {
				action1, nonMember := &Blank{"member"}, &Blank{"non-member"}
				pipeline := NewPipeline("Pipeline", action1)

				pipeline.SetRunPlan(action1, SuccessOnlyPlan(nonMember))
			},
			panics:  true,
			message: "setting plan from `member` directing `success` to non-member `non-member`",
		},
		"set action plan with self loop": {
			code: func() {
				action1 := &Blank{"member"}
				pipeline := NewPipeline("Pipeline", action1)

				pipeline.SetRunPlan(action1, SuccessOnlyPlan(action1))
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

func TestRecoverOnRun(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)

	ctx := context.Background()
	type testCase struct {
		closure func(t *testing.T) func()
	}

	testCases := map[string]testCase{
		"recover from action": {
			func(t *testing.T) func() {
				return func() {
					pipeline := NewPipeline("Calculation", &Divide{})
					_, err := pipeline.Run(ctx, 0)

					assert.Error(t, err)
					assert.Contains(t, err.Error(), "divide by zero")
				}
			}},
		"skip actions after recover": {
			func(t *testing.T) func() {
				return func() {
					pipeline := NewPipeline("Calculation", &Divide{}, &SetTen{})
					output, err := pipeline.Run(ctx, 0)

					assert.Error(t, err)
					assert.Contains(t, err.Error(), "divide by zero")
					assert.NotEqual(t, 10, output)
				}
			}},
		"internal pipeline panic recovers": {
			func(t *testing.T) func() {
				return func() {
					subPipeline := NewPipeline("SubPipeline", &Divide{}, &SetTen{})
					pipeline := NewPipeline(
						"SuperPipeline",
						&SetTen{name: "setTen"},
						&Divide{name: "divide"},
						subPipeline,
						&SetTen{name: "runAfterPanic"},
					)
					output, err := pipeline.Run(ctx, 0)

					assert.Error(t, err)
					assert.Contains(t, err.Error(), "divide by zero")
					assert.NotEqual(t, 10, output)
				}
			}},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			code := tc.closure(t)
			assert.NotPanics(t, code)
		})
	}
}

type Blank struct{ name string }

func (b Blank) Name() string { return b.name }
func (Blank) Run(_ context.Context, _ string) (string, error) {
	return "", nil
}

type Divide struct{ name string }

func (Divide) Name() string { return "Divide" }
func (Divide) Run(_ context.Context, input int) (output int, err error) {
	// panics when input is zero
	return 1 / input, nil
}

type SetTen struct{ name string }

func (SetTen) Name() string { return "SetTen" }
func (SetTen) Run(_ context.Context, _ int) (output int, err error) {
	return 10, nil
}
