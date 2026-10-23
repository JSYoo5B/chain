package numbers

import (
	"context"
	"errors"

	"github.com/JSYoo5B/chain"
)

func increment() chain.Action[int] {
	return chain.NewSimpleAction("increment", func(_ context.Context, input int) (int, error) {
		return input + 1, nil
	})
}

func double() chain.Action[int] {
	return chain.NewSimpleAction("double", func(_ context.Context, input int) (int, error) {
		return input * 2, nil
	})
}

func addTen() chain.Action[int] {
	return chain.NewSimpleAction("add ten", func(_ context.Context, input int) (int, error) {
		return input + 10, nil
	})
}

func linear() *chain.Workflow[int] {
	return chain.NewWorkflow("numbers", increment(), double())
}

func runFromDouble(ctx context.Context, input int) (int, error) {
	first, second := increment(), double()
	workflow := chain.NewWorkflow("numbers", first, second)
	return workflow.RunAt(second, ctx, input)
}

func skipDouble() *chain.Workflow[int] {
	first, second, last := increment(), double(), addTen()
	workflow := chain.NewWorkflow("skip double", first, second, last)
	workflow.SetRunPlan(first, chain.RunPlan[int]{
		chain.Success: last,
		chain.Failure: chain.Terminate[int](),
		chain.Abort:   chain.Terminate[int](),
	})
	return workflow
}

func stopAtDouble() *chain.Workflow[int] {
	first, second, last := increment(), double(), addTen()
	workflow := chain.NewWorkflow("stop at double", first, second, last)
	workflow.SetRunPlan(first, chain.DefaultPlan(second, last))
	workflow.SetRunPlan(second, chain.TerminationPlan[int]())
	return workflow
}

var errAddition = errors.New("addition failed")

func handleFailure() *chain.Workflow[int] {
	first := chain.NewSimpleAction("failed increment", func(_ context.Context, input int) (int, error) {
		return input + 1, errAddition
	})
	second, last := double(), addTen()
	workflow := chain.NewWorkflow("handle failure", first, second, last)
	workflow.SetRunPlan(first, chain.DefaultPlan(second, last))
	workflow.SetRunPlan(second, chain.TerminationPlan[int]())
	return workflow
}

// Select Abort explicitly for negative results, so this example needs no panic.
func resetNegative() *chain.Workflow[int] {
	first := chain.AsBranchAction(increment(), func(_ context.Context, output int) (string, error) {
		if output < 0 {
			return chain.Abort, nil
		}
		return chain.Success, nil
	})
	second, onFailure := double(), addTen()
	reset := chain.NewSimpleAction("reset", func(_ context.Context, _ int) (int, error) {
		return 0, nil
	})
	workflow := chain.NewWorkflow("reset negatives", first, second, onFailure, reset)
	workflow.SetRunPlan(first, chain.DefaultPlanWithAbort(second, onFailure, reset))
	workflow.SetRunPlan(second, chain.TerminationPlan[int]())
	workflow.SetRunPlan(onFailure, chain.TerminationPlan[int]())
	return workflow
}

// Construct an invalid graph for the validation example; do not execute it.
func cycle() *chain.Workflow[int] {
	first, second := increment(), double()
	workflow := chain.NewWorkflow("cycle", first, second)
	workflow.SetRunPlan(second, chain.SuccessOnlyPlan(first))
	return workflow
}
