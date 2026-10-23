package errorcontrol

import (
	"context"
	"fmt"
)

func Example_retryWithRollback() {
	action, counts := retryingIncrement(true)
	output, err := action.Run(context.Background(), 0)
	fmt.Println(output, counts.Calls, counts.Rollbacks, err)
	// Output: 1 3 2 <nil>
}

func Example_retryWithoutRollback() {
	action, counts := retryingIncrement(false)
	output, err := action.Run(context.Background(), 0)
	fmt.Println(output, counts.Calls, counts.Rollbacks, err)
	// Output: 3 3 0 <nil>
}

func Example_optionalIncrement() {
	action, fallback := optionalIncrement()
	output, err := action.Run(context.Background(), 5)
	fmt.Println("fallback:", fallback.Input, fallback.Err)
	fmt.Println("result:", output, err)
	// Output:
	// fallback: 5 optional step failed
	// result: 6 <nil>
}
