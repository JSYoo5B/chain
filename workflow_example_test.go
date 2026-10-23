package chain_test

import (
	"context"
	"errors"
	"fmt"

	"github.com/JSYoo5B/chain"
)

// NewWorkflow connects successful actions in constructor order.
func ExampleNewWorkflow() {
	increment := chain.NewSimpleAction("increment", func(_ context.Context, input int) (int, error) {
		return input + 1, nil
	})
	double := chain.NewSimpleAction("double", func(_ context.Context, input int) (int, error) {
		return input * 2, nil
	})
	workflow := chain.NewWorkflow("numbers", increment, double)
	output, err := workflow.Run(context.Background(), 2)
	fmt.Println(workflow.Name(), output, err)
	// Output: numbers 6 <nil>
}

// RunAt starts at the chosen member, skipping earlier actions.
func ExampleWorkflow_RunAt() {
	increment := chain.NewSimpleAction("increment", func(_ context.Context, input int) (int, error) {
		return input + 1, nil
	})
	double := chain.NewSimpleAction("double", func(_ context.Context, input int) (int, error) {
		return input * 2, nil
	})
	workflow := chain.NewWorkflow("numbers", increment, double)
	output, err := workflow.RunAt(double, context.Background(), 3)
	fmt.Println(output, err)
	// Output: 6 <nil>
}

// An explicit plan can skip the next declared member.
func ExampleWorkflow_SetRunPlan() {
	increment := chain.NewSimpleAction("increment", func(_ context.Context, input int) (int, error) {
		return input + 1, nil
	})
	double := chain.NewSimpleAction("double", func(_ context.Context, input int) (int, error) {
		return input * 2, nil
	})
	addTen := chain.NewSimpleAction("add ten", func(_ context.Context, input int) (int, error) {
		return input + 10, nil
	})
	workflow := chain.NewWorkflow("skip double", increment, double, addTen)
	workflow.SetRunPlan(increment, chain.RunPlan[int]{chain.Success: addTen})
	output, err := workflow.Run(context.Background(), 0)
	fmt.Println(output, err)
	// Output: 11 <nil>
}

// ValidateGraph accepts a linear graph and detects a subsequently added cycle.
func ExampleWorkflow_ValidateGraph() {
	identity := func(_ context.Context, input int) (int, error) { return input, nil }
	first := chain.NewSimpleAction("first", identity)
	second := chain.NewSimpleAction("second", identity)
	workflow := chain.NewWorkflow("numbers", first, second)
	fmt.Println("valid:", workflow.ValidateGraph() == nil)
	workflow.SetRunPlan(second, chain.SuccessOnlyPlan(first))
	fmt.Println("cycle:", workflow.ValidateGraph() != nil)
	// Output:
	// valid: true
	// cycle: true
}

// Only the success direction has a next action.
func ExampleSuccessOnlyPlan() {
	double := chain.NewSimpleAction("double", func(_ context.Context, input int) (int, error) {
		return input * 2, nil
	})
	plan := chain.SuccessOnlyPlan(double)
	fmt.Println(plan[chain.Success].Name())
	fmt.Println(plan[chain.Failure] == nil, plan[chain.Abort] == nil)
	// Output:
	// double
	// true true
}

// A failure handler receives the failed output; handling preserves the error.
func ExampleDefaultPlan() {
	failure := errors.New("add one failed")
	step := chain.NewSimpleAction("add one", func(_ context.Context, input int) (int, error) {
		return input + 1, failure
	})
	double := chain.NewSimpleAction("double", func(_ context.Context, input int) (int, error) {
		return input * 2, nil
	})
	addTen := chain.NewSimpleAction("add ten", func(_ context.Context, input int) (int, error) {
		return input + 10, nil
	})
	workflow := chain.NewWorkflow("numbers", step, double, addTen)
	workflow.SetRunPlan(step, chain.DefaultPlan(double, addTen))
	workflow.SetRunPlan(double, chain.TerminationPlan[int]())
	output, err := workflow.Run(context.Background(), 0)
	fmt.Println(output, errors.Is(err, failure))
	// Output: 11 true
}

// Success, failure, and abort can each select a different action.
func ExampleDefaultPlanWithAbort() {
	identity := func(_ context.Context, input int) (int, error) { return input, nil }
	next := chain.NewSimpleAction("next", identity)
	onFailure := chain.NewSimpleAction("on failure", identity)
	onAbort := chain.NewSimpleAction("on abort", identity)
	plan := chain.DefaultPlanWithAbort(next, onFailure, onAbort)
	fmt.Println(plan[chain.Success].Name())
	fmt.Println(plan[chain.Failure].Name())
	fmt.Println(plan[chain.Abort].Name())
	// Output:
	// next
	// on failure
	// on abort
}

// A leaf terminates instead of falling through to the next declared member.
func ExampleTerminationPlan() {
	increment := chain.NewSimpleAction("increment", func(_ context.Context, input int) (int, error) {
		return input + 1, nil
	})
	double := chain.NewSimpleAction("double", func(_ context.Context, input int) (int, error) {
		return input * 2, nil
	})
	addTen := chain.NewSimpleAction("add ten", func(_ context.Context, input int) (int, error) {
		return input + 10, nil
	})
	workflow := chain.NewWorkflow("stop at double", increment, double, addTen)
	workflow.SetRunPlan(increment, chain.DefaultPlan(double, addTen))
	workflow.SetRunPlan(double, chain.TerminationPlan[int]())
	output, err := workflow.Run(context.Background(), 3)
	fmt.Println(output, err)
	// Output: 8 <nil>
}

// Terminate is a destination in a plan; it is not a workflow member.
func ExampleTerminate() {
	double := chain.NewSimpleAction("double", func(_ context.Context, input int) (int, error) {
		return input * 2, nil
	})
	plan := chain.RunPlan[int]{
		chain.Success: double,
		chain.Failure: chain.Terminate[int](),
		chain.Abort:   chain.Terminate[int](),
	}
	fmt.Println(plan[chain.Success].Name(), plan[chain.Failure] == nil, plan[chain.Abort] == nil)
	// Output: double true true
}
