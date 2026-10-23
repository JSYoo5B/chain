package numbers

import (
	"context"
	"errors"
	"fmt"
)

func Example_double() {
	action := double()
	output, err := action.Run(context.Background(), 3)
	fmt.Println(action.Name(), output, err)
	// Output: double 6 <nil>
}

func Example_linear() {
	workflow := linear()
	output, err := workflow.Run(context.Background(), 2)
	fmt.Println(workflow.Name(), output, err)
	// Output: numbers 6 <nil>
}

func Example_startAtDouble() {
	output, err := runFromDouble(context.Background(), 3)
	fmt.Println(output, err)
	// Output: 6 <nil>
}

func Example_skipDouble() {
	output, err := skipDouble().Run(context.Background(), 0)
	fmt.Println(output, err)
	// Output: 11 <nil>
}

func Example_stopAtDouble() {
	output, err := stopAtDouble().Run(context.Background(), 3)
	fmt.Println(output, err)
	// Output: 8 <nil>
}

func Example_handleFailure() {
	output, err := handleFailure().Run(context.Background(), 0)
	fmt.Println(output, errors.Is(err, errAddition))
	// Output: 11 true
}

func Example_resetNegative() {
	for _, input := range []int{-3, 3} {
		output, err := resetNegative().Run(context.Background(), input)
		fmt.Println(input, output, err)
	}
	// Output:
	// -3 0 <nil>
	// 3 8 <nil>
}

func Example_graphValidation() {
	fmt.Println("valid:", linear().ValidateGraph() == nil)
	fmt.Println("cycle:", cycle().ValidateGraph() != nil)
	// Output:
	// valid: true
	// cycle: true
}
