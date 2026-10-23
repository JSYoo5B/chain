package branch

import (
	"context"
	"fmt"
)

func Example_checkSign() {
	action := checkSign()
	ctx := context.Background()
	output, err := action.Run(ctx, -2)
	direction, branchErr := action.NextDirection(ctx, output)
	fmt.Println(output, direction, err, branchErr)
	// Output: -2 negative <nil> <nil>
}

func Example_incrementAndCheckParity() {
	action := incrementAndCheckParity()
	ctx := context.Background()
	output, err := action.Run(ctx, 3)
	direction, branchErr := action.NextDirection(ctx, output)
	fmt.Println(action.Name(), action.Directions())
	fmt.Println(output, direction, err, branchErr)
	// Output:
	// Inc [even odd]
	// 4 even <nil> <nil>
}
