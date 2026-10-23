package chain_test

import (
	"context"
	"fmt"

	"github.com/JSYoo5B/chain"
)

// Adapt an integer action to one field, retaining the struct's other fields.
func ExampleAdaptAction() {
	type number struct {
		Value int
		Label string
	}
	increment := chain.NewSimpleAction("increment", func(_ context.Context, input int) (int, error) {
		return input + 1, nil
	})
	adapted := chain.AdaptAction(increment,
		func(input number) int { return input.Value },
		func(input number, value int) number {
			input.Value = value
			return input
		})
	output, err := adapted.Run(context.Background(), number{Value: 2, Label: "keep"})
	fmt.Println(output, err)
	// Output: {3 keep} <nil>
}
