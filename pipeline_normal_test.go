package railway

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSimplePipeline(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)

	collatz := NewCollatz("SingleCollatz")
	ctx := context.Background()
	t.Run("Odd input collatz", func(t *testing.T) {
		output, err := collatz.Run(ctx, 5)

		assert.NoError(t, err)
		assert.Equal(t, 16, output)
	})
	t.Run("Even input collatz", func(t *testing.T) {
		output, err := collatz.Run(ctx, 16)

		assert.NoError(t, err)
		assert.Equal(t, 8, output)
	})
	t.Run("RunAt onOdd with even input", func(t *testing.T) {
		output, err := collatz.RunAt(collatz.OnOdd, ctx, 2)

		assert.NoError(t, err)
		assert.Equal(t, 7, output)
	})
	t.Run("RunAt onEven with odd input", func(t *testing.T) {
		_, err := collatz.RunAt(collatz.OnEven, ctx, 5)

		assert.Error(t, err)
	})
}

func TestMultiPipeline(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)

	multiCollatz := NewPipeline(
		"FiveCollatzTrial",
		NewCollatz("collatz1"),
		NewPipeline(
			"InnerCollatz",
			NewCollatz("collatz2"),
			NewCollatz("collatz3"),
		),
		NewCollatz("collatz4"),
		NewCollatz("collatz5"),
	)
	ctx := context.Background()
	t.Run("Start with 5", func(t *testing.T) {
		// 5 -> 16 -> 8 -> 4 -> 2 -> 1
		output, err := multiCollatz.Run(ctx, 5)

		assert.NoError(t, err)
		assert.Equal(t, 1, output)
	})
	t.Run("Start with 10", func(t *testing.T) {
		// 10 -> 5 -> 16 -> 8 -> 4 -> 2
		output, err := multiCollatz.Run(ctx, 10)

		assert.NoError(t, err)
		assert.Equal(t, 2, output)
	})
	t.Run("Start with 1", func(t *testing.T) {
		// 1 -> 1 -> 1 -> 1 -> 1 -> 1
		output, err := multiCollatz.Run(ctx, 1)

		assert.NoError(t, err)
		assert.Equal(t, 1, output)
	})
	t.Run("Start with 0", func(t *testing.T) {
		// 0 -> Finish (terminated by error)
		_, err := multiCollatz.Run(ctx, 0)

		assert.Error(t, err)
	})
}

func TestErrorPropagation(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)

	t.Run("simple pipeline with single error", func(t *testing.T) {
		action1 := NewCollatz("action1")
		action2 := &ErrorMaker{message: "error2"}
		action3 := NewCollatz("action3")
		pipeline := NewPipeline("Pipeline", action1, action2, action3)
		pipeline.SetRunPlan(action1, SuccessOnlyPlan(action2))
		pipeline.SetRunPlan(action2, DefaultPlan(action3, action3))
		pipeline.SetRunPlan(action3, TerminationPlan[int]())

		// 5 -> 16 -> 16(error2) -> 8
		output, err := pipeline.Run(context.Background(), 5)

		assert.Error(t, err)
		assert.Equal(t, err.Error(), "error2")
		assert.Equal(t, 8, output)
	})

	t.Run("simple pipeline with multiple errors", func(t *testing.T) {
		action1 := NewCollatz("action1")
		action2 := &ErrorMaker{message: "error2"}
		action3 := &ErrorMaker{message: "error3"}
		pipeline := NewPipeline("Pipeline", action1, action2, action3)
		pipeline.SetRunPlan(action1, SuccessOnlyPlan(action2))
		pipeline.SetRunPlan(action2, DefaultPlan(action3, action3))
		pipeline.SetRunPlan(action3, TerminationPlan[int]())

		// 5 -> 16 -> 16(error2) -> 16(error3)
		output, err := pipeline.Run(context.Background(), 5)

		assert.Error(t, err)
		assert.Equal(t, err.Error(), "error3")
		assert.Equal(t, 16, output)
	})

	t.Run("simple pipeline with panic", func(t *testing.T) {
		action1 := NewCollatz("action1")
		action2 := &PanicMaker{message: "panic2"}
		action3 := NewCollatz("action3")
		pipeline := NewPipeline("Pipeline", action1, action2, action3)
		pipeline.SetRunPlan(action1, SuccessOnlyPlan(action2))
		pipeline.SetRunPlan(action2, DefaultPlan(action3, action3))
		pipeline.SetRunPlan(action3, TerminationPlan[int]())

		// 5 -> 16 -> 16(panic2) -> abort
		output, err := pipeline.Run(context.Background(), 5)

		assert.Error(t, err)
		assert.Equal(t, err.Error(), "panic2")
		assert.Equal(t, 16, output)
	})

	t.Run("simple pipeline with panic and abort planning", func(t *testing.T) {
		action1 := NewCollatz("action1")
		action2 := &PanicMaker{message: "panic2"}
		action3 := NewCollatz("action3")
		pipeline := NewPipeline("Pipeline", action1, action2, action3)
		pipeline.SetRunPlan(action1, SuccessOnlyPlan(action2))
		pipeline.SetRunPlan(action2, DefaultPlanWithAbort(action3, action3, action3))
		pipeline.SetRunPlan(action3, TerminationPlan[int]())

		// 5 -> 16 -> 16(panic2) -(abort)-> 8
		output, err := pipeline.Run(context.Background(), 5)

		assert.Error(t, err)
		assert.Equal(t, err.Error(), "panic2")
		assert.Equal(t, 8, output)
	})

	t.Run("multi pipeline with internal error", func(t *testing.T) {
		action1 := NewCollatz("action1")
		action2 := &ErrorMaker{message: "error2"}
		action3 := NewCollatz("action3")
		subPipeline := NewPipeline("SubPipeline", action1, action2, action3)
		subPipeline.SetRunPlan(action1, SuccessOnlyPlan(action2))
		subPipeline.SetRunPlan(action2, DefaultPlan(action3, action3))
		subPipeline.SetRunPlan(action3, TerminationPlan[int]())
		action0 := NewCollatz("action0")
		action4 := NewCollatz("action4")
		pipeline := NewPipeline("Pipeline", action0, subPipeline, action4)
		pipeline.SetRunPlan(action0, SuccessOnlyPlan(subPipeline))
		pipeline.SetRunPlan(subPipeline, DefaultPlan(action4, action4))
		pipeline.SetRunPlan(action4, TerminationPlan[int]())

		// 5 -> [16 -> 8 -> 8(error2) -> 4] -> 2
		output, err := pipeline.Run(context.Background(), 5)

		assert.Error(t, err)
		assert.Equal(t, err.Error(), "error2")
		assert.Equal(t, 2, output)
	})
}

type Collatz struct {
	*Pipeline[int]
	CheckNext Action[int]
	OnEven    Action[int]
	OnOdd     Action[int]
}

func NewCollatz(name string) *Collatz {
	checkNext, onEven, onOdd := &CheckNext{}, &OnEven{}, &OnOdd{}
	pipeline := NewPipeline(
		name,
		checkNext,
		onEven,
		onOdd,
	)

	terminate := Terminate[int]()
	noActionPlan := TerminationPlan[int]()
	pipeline.SetRunPlan(checkNext, ActionPlan[int]{
		"even": onEven,
		"odd":  onOdd,
	})
	pipeline.SetRunPlan(onEven, ActionPlan[int]{
		Success: terminate,
		// Skip Error, automatically configured as terminate
	})
	pipeline.SetRunPlan(onOdd, noActionPlan)

	return &Collatz{
		Pipeline:  pipeline,
		CheckNext: checkNext,
		OnEven:    onEven,
		OnOdd:     onOdd,
	}
}

type CheckNext struct{}

func (CheckNext) Name() string         { return "CheckNext" }
func (CheckNext) Directions() []string { return []string{"even", "odd"} }
func (CheckNext) Run(_ context.Context, input int) (output int, err error) {
	return input, nil
}
func (CheckNext) NextDirection(_ context.Context, output int) (direction string, err error) {
	if output == 1 {
		return Abort, nil
	} else if output == 0 {
		return Error, fmt.Errorf("cannot try collatz with 0")
	} else if output%2 == 0 {
		return "even", nil
	} else {
		return "odd", nil
	}
}

type OnEven struct{}

func (OnEven) Name() string { return "OnEven" }
func (OnEven) Run(_ context.Context, input int) (output int, err error) {
	if input%2 != 0 {
		// direction is intended bug. on running pipeline, it should be changed as Error
		return input, fmt.Errorf("input is not even")
	}
	return input / 2, nil
}

type OnOdd struct{}

func (OnOdd) Name() string { return "OnOdd" }
func (OnOdd) Run(_ context.Context, input int) (output int, err error) {
	return 3*input + 1, nil
}

type ErrorMaker struct{ message string }

func (e ErrorMaker) Name() string { return e.message }
func (e ErrorMaker) Run(_ context.Context, input int) (output int, err error) {
	return input, errors.New(e.message)
}

type PanicMaker struct{ message string }

func (p PanicMaker) Name() string { return p.message }
func (p PanicMaker) Run(_ context.Context, _ int) (output int, err error) {
	panic(errors.New(p.message))
}
