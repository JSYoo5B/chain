package collection

import (
	"context"
	"fmt"
)

func Example_sequenceNumbers() {
	output, err := sequenceNumbers().Run(context.Background(), []int{1, 2, 3})
	fmt.Println(output, err)
	// Output: [2 4 6] <nil>
}

func Example_stopOnError() {
	for _, stopOnError := range []bool{true, false} {
		output, err := sequenceWithError(stopOnError).Run(context.Background(), []int{1, 2, 3})
		fmt.Println(stopOnError, output, err != nil)
	}
	// Output:
	// true [2 2 3] true
	// false [2 2 6] true
}

func Example_sequenceWords() {
	output, err := sequenceWords().Run(context.Background(), map[string]string{"a": "go", "b": "chain"})
	fmt.Println(output["a"], output["b"], err)
	// Output: GO CHAIN <nil>
}

func Example_parallelNumbers() {
	output, err := parallelNumbers().Run(context.Background(), []int{1, 2, 3})
	fmt.Println(output, err)
	// Output: [2 4 6] <nil>
}

func Example_parallelNumberMap() {
	output, err := parallelNumberMap().Run(context.Background(), map[string]int{"a": 1, "b": 2})
	fmt.Println(output["a"], output["b"], err)
	// Output: 2 4 <nil>
}
