package chain_test

import (
	"context"
	"errors"
	"fmt"

	"github.com/JSYoo5B/chain"
)

// The three-attempt budget permits two rollbacks before the successful attempt.
func ExampleAsRetryableAction() {
	attempts := 0
	step := chain.NewSimpleAction("add one", func(_ context.Context, input int) (int, error) {
		attempts++
		if attempts < 3 {
			return input + 1, errors.New("not yet")
		}
		return input + 1, nil
	})
	undo := chain.NewSimpleAction("subtract one", func(_ context.Context, input int) (int, error) {
		fmt.Println("rollback:", input)
		return input - 1, nil
	})
	retry := chain.AsRetryableAction("retry", step, undo, 3)
	output, err := retry.Run(context.Background(), 0)
	fmt.Println("result:", output, "attempts:", attempts, err)
	// Output:
	// rollback: 1
	// rollback: 1
	// result: 1 attempts: 3 <nil>
}

// Without rollback, the next attempt receives the failed attempt's output.
func ExampleSkipRollback() {
	attempts := 0
	step := chain.NewSimpleAction("add one", func(_ context.Context, input int) (int, error) {
		attempts++
		if attempts == 1 {
			return input + 1, errors.New("try again")
		}
		return input + 1, nil
	})
	retry := chain.AsRetryableAction("retry", step, chain.SkipRollback[int](), 2)
	output, err := retry.Run(context.Background(), 0)
	fmt.Println(output, attempts, err)
	// Output: 2 2 <nil>
}

// The fallback sees the original input; the returned value keeps the action's output.
func ExampleAsBestEffortAction() {
	step := chain.NewSimpleAction("add one", func(_ context.Context, input int) (int, error) {
		return input + 1, errors.New("optional step failed")
	})
	optional := chain.AsBestEffortAction(step, func(_ context.Context, input int, err error) {
		fmt.Println("fallback:", input, err)
	})
	output, err := optional.Run(context.Background(), 5)
	fmt.Println(output, err)
	// Output:
	// fallback: 5 optional step failed
	// 6 <nil>
}
