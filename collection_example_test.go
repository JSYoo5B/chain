package chain_test

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/JSYoo5B/chain"
)

// Sequential slice processing keeps index order.
func ExampleAsSequenceSliceAction() {
	double := chain.NewSimpleAction("double", func(_ context.Context, input int) (int, error) {
		return input * 2, nil
	})
	each := chain.AsSequenceSliceAction("double each", double, true)
	output, err := each.Run(context.Background(), []int{1, 2, 3})
	fmt.Println(output, err)
	// Output: [2 4 6] <nil>
}

// Stopping on error leaves later values unchanged; continuing processes them.
func ExampleAsSequenceSliceAction_stopOnError() {
	double := chain.NewSimpleAction("double except two", func(_ context.Context, input int) (int, error) {
		if input == 2 {
			return input, errors.New("two is unchanged")
		}
		return input * 2, nil
	})
	for _, stopOnError := range []bool{true, false} {
		each := chain.AsSequenceSliceAction("double each", double, stopOnError)
		output, err := each.Run(context.Background(), []int{1, 2, 3})
		fmt.Println(stopOnError, output, err != nil)
	}
	// Output:
	// true [2 2 3] true
	// false [2 2 6] true
}

// Map keys stay the same while the values are transformed.
func ExampleAsSequenceMapAction() {
	upper := chain.NewSimpleAction("uppercase", func(_ context.Context, input string) (string, error) {
		return strings.ToUpper(input), nil
	})
	each := chain.AsSequenceMapAction[string]("uppercase each", upper)
	output, err := each.Run(context.Background(), map[string]string{"a": "go", "b": "chain"})
	// Read keys explicitly rather than depending on map iteration order.
	fmt.Println(output["a"], output["b"], err)
	// Output: GO CHAIN <nil>
}

// A pure function can run concurrently while output positions remain stable.
func ExampleAsParallelSliceAction() {
	double := chain.NewSimpleAction("double", func(_ context.Context, input int) (int, error) {
		return input * 2, nil
	})
	each := chain.AsParallelSliceAction("double each", double)
	output, err := each.Run(context.Background(), []int{1, 2, 3})
	fmt.Println(output, err)
	// Output: [2 4 6] <nil>
}

// Concurrent map processing preserves each value's key.
func ExampleAsParallelMapAction() {
	double := chain.NewSimpleAction("double", func(_ context.Context, input int) (int, error) {
		return input * 2, nil
	})
	each := chain.AsParallelMapAction[string]("double each", double)
	output, err := each.Run(context.Background(), map[string]int{"a": 1, "b": 2})
	fmt.Println(output["a"], output["b"], err)
	// Output: 2 4 <nil>
}
