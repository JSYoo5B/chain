package dag

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
		output, direction, err := collatz.Run(ctx, 5)

		assert.NoError(t, err)
		assert.Equal(t, 16, output)
		assert.Equal(t, Success, direction)
	})
	t.Run("Even input collatz", func(t *testing.T) {
		output, direction, err := collatz.Run(ctx, 16)

		assert.NoError(t, err)
		assert.Equal(t, 8, output)
		assert.Equal(t, Success, direction)
	})
	t.Run("RunAt onOdd with even input", func(t *testing.T) {
		output, direction, err := collatz.RunAt(collatz.OnOdd, ctx, 2)

		assert.NoError(t, err)
		assert.Equal(t, 7, output)
		assert.Equal(t, Success, direction)
	})
	t.Run("RunAt onEven with odd input", func(t *testing.T) {
		_, direction, err := collatz.RunAt(collatz.OnEven, ctx, 5)

		assert.Error(t, err)
		assert.Equal(t, Error, direction)
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
		output, direction, err := multiCollatz.Run(ctx, 5)

		assert.NoError(t, err)
		assert.Equal(t, 1, output)
		assert.Equal(t, Success, direction)
	})
	t.Run("Start with 10", func(t *testing.T) {
		// 10 -> 5 -> 16 -> 8 -> 4 -> 2
		output, direction, err := multiCollatz.Run(ctx, 10)

		assert.NoError(t, err)
		assert.Equal(t, 2, output)
		assert.Equal(t, Success, direction)
	})
	t.Run("Start with 1", func(t *testing.T) {
		// 1 -> 1 -> 1 -> 1 -> 1 -> 1
		output, direction, err := multiCollatz.Run(ctx, 1)

		assert.NoError(t, err)
		assert.Equal(t, 1, output)
		assert.Equal(t, Abort, direction)
	})
	t.Run("Start with 0", func(t *testing.T) {
		// 0 -> Finish (terminated by error)
		_, direction, err := multiCollatz.Run(ctx, 0)

		assert.Error(t, err)
		assert.Equal(t, Error, direction)
	})
}

func TestErrorPropagation(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)

	t.Run("simple pipeline with single error", func(t *testing.T) {
		action1 := NewCollatz("action1")
		action2 := &ErrorMaker{message: "error2"}
		action3 := NewCollatz("action3")
		pipeline := NewPipeline("Pipeline", action1, action2, action3)
		pipeline.SetRunPlan(action1, ActionPlan[int]{Success: action2})
		pipeline.SetRunPlan(action2, ActionPlan[int]{Success: action3, Error: action3})
		pipeline.SetRunPlan(action3, TerminationPlan[int]())

		// 5 -> 16 -> 16(error2) -> 8
		output, direction, err := pipeline.Run(context.Background(), 5)

		assert.Error(t, err)
		assert.Equal(t, err.Error(), "error2")
		assert.Equal(t, Error, direction)
		assert.Equal(t, 8, output)
	})

	t.Run("simple pipeline with multiple errors", func(t *testing.T) {
		action1 := NewCollatz("action1")
		action2 := &ErrorMaker{message: "error2"}
		action3 := &ErrorMaker{message: "error3"}
		pipeline := NewPipeline("Pipeline", action1, action2, action3)
		pipeline.SetRunPlan(action1, ActionPlan[int]{Success: action2})
		pipeline.SetRunPlan(action2, ActionPlan[int]{Success: action3, Error: action3})
		pipeline.SetRunPlan(action3, TerminationPlan[int]())

		// 5 -> 16 -> 16(error2) -> 16(error3)
		output, direction, err := pipeline.Run(context.Background(), 5)

		assert.Error(t, err)
		assert.Equal(t, err.Error(), "error3")
		assert.Equal(t, Error, direction)
		assert.Equal(t, 16, output)
	})

	t.Run("simple pipeline with panic", func(t *testing.T) {
		action1 := NewCollatz("action1")
		action2 := &PanicMaker{message: "panic2"}
		action3 := NewCollatz("action3")
		pipeline := NewPipeline("Pipeline", action1, action2, action3)
		pipeline.SetRunPlan(action1, ActionPlan[int]{Success: action2})
		pipeline.SetRunPlan(action2, ActionPlan[int]{Success: action3, Error: action3})
		pipeline.SetRunPlan(action3, TerminationPlan[int]())

		// 5 -> 16 -> 16(panic2) -> abort
		output, direction, err := pipeline.Run(context.Background(), 5)

		assert.Error(t, err)
		assert.Equal(t, err.Error(), "panic2")
		assert.Equal(t, Abort, direction)
		assert.Equal(t, 16, output)
	})

	t.Run("multi pipeline with internal error", func(t *testing.T) {
		action1 := NewCollatz("action1")
		action2 := &ErrorMaker{message: "error2"}
		action3 := NewCollatz("action3")
		subPipeline := NewPipeline("SubPipeline", action1, action2, action3)
		subPipeline.SetRunPlan(action1, ActionPlan[int]{Success: action2})
		subPipeline.SetRunPlan(action2, ActionPlan[int]{Success: action3, Error: action3})
		subPipeline.SetRunPlan(action3, TerminationPlan[int]())
		action0 := NewCollatz("action0")
		action4 := NewCollatz("action4")
		pipeline := NewPipeline("Pipeline", action0, subPipeline, action4)
		pipeline.SetRunPlan(action0, ActionPlan[int]{Success: subPipeline})
		pipeline.SetRunPlan(subPipeline, ActionPlan[int]{Success: action4, Error: action4})
		pipeline.SetRunPlan(action4, TerminationPlan[int]())

		// 5 -> [16 -> 8 -> 8(error2) -> 4] -> 2
		output, direction, err := pipeline.Run(context.Background(), 5)

		assert.Error(t, err)
		assert.Equal(t, err.Error(), "error2")
		assert.Equal(t, Error, direction)
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

func (CheckNext) Name() string { return "CheckNext" }
func (CheckNext) Run(_ context.Context, input int) (output int, direction string, err error) {
	if input == 1 {
		return input, Abort, nil
	} else if input == 0 {
		return input, Error, fmt.Errorf("cannot try collatz with 0")
	} else if input%2 == 0 {
		return input, "even", nil
	} else {
		return input, "odd", nil
	}
}

type OnEven struct{}

func (OnEven) Name() string { return "OnEven" }
func (OnEven) Run(_ context.Context, input int) (output int, direction string, err error) {
	if input%2 != 0 {
		// direction is intended bug. on running pipeline, it should be changed as Error
		return input, Success, fmt.Errorf("input is not even")
	}
	return input / 2, Success, nil
}

type OnOdd struct{}

func (OnOdd) Name() string { return "OnOdd" }
func (OnOdd) Run(_ context.Context, input int) (output int, direction string, err error) {
	return 3*input + 1, Success, nil
}

type ErrorMaker struct{ message string }

func (e ErrorMaker) Name() string { return e.message }
func (e ErrorMaker) Run(_ context.Context, input int) (output int, direction string, err error) {
	return input, Error, errors.New(e.message)
}

type PanicMaker struct{ message string }

func (p PanicMaker) Name() string { return p.message }
func (p PanicMaker) Run(_ context.Context, _ int) (output int, direction string, err error) {
	panic(errors.New(p.message))
}
