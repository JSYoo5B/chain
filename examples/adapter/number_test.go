package adapter

import (
	"context"
	"fmt"
)

func Example_incrementNumber() {
	output, err := incrementNumber().Run(context.Background(), number{Value: 2, Label: "keep"})
	fmt.Println(output, err)
	// Output: {3 keep} <nil>
}
