package chain_test

import (
	"context"
	"fmt"

	"github.com/JSYoo5B/chain"
)

// A function becomes an action with a name and a Run method.
func ExampleNewSimpleAction() {
	double := chain.NewSimpleAction("double", func(_ context.Context, input int) (int, error) {
		return input * 2, nil
	})
	output, err := double.Run(context.Background(), 3)
	fmt.Println(double.Name(), output, err)
	// Output: double 6 <nil>
}

// A nil run function passes the input through unchanged before choosing a branch.
func ExampleNewSimpleBranchAction() {
	sign := chain.NewSimpleBranchAction[int]("sign", nil, []string{"negative", "nonnegative"},
		func(_ context.Context, output int) (string, error) {
			if output < 0 {
				return "negative", nil
			}
			return "nonnegative", nil
		})
	ctx := context.Background()
	output, err := sign.Run(ctx, -2)
	direction, branchErr := sign.NextDirection(ctx, output)
	fmt.Println(output, direction, err, branchErr)
	// Output: -2 negative <nil> <nil>
}

// The wrapped action runs first; the selector receives its output, not its input.
func ExampleAsBranchAction() {
	increment := chain.NewSimpleAction("increment", func(_ context.Context, input int) (int, error) {
		return input + 1, nil
	})
	parity := chain.AsBranchAction(increment, func(_ context.Context, output int) (string, error) {
		if output%2 == 0 {
			return "even", nil
		}
		return "odd", nil
	}, "even", "odd")
	ctx := context.Background()
	output, err := parity.Run(ctx, 3)
	direction, branchErr := parity.NextDirection(ctx, output)
	fmt.Println(parity.Name(), parity.Directions())
	fmt.Println(output, direction, err, branchErr)
	// Output:
	// increment [even odd]
	// 4 even <nil> <nil>
}
