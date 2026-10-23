package adapter

import (
	"context"

	"github.com/JSYoo5B/chain"
)

type number struct {
	Value int
	Label string
}

func incrementNumber() chain.Action[number] {
	increment := chain.NewSimpleAction("increment", func(_ context.Context, input int) (int, error) {
		return input + 1, nil
	})
	return chain.AdaptAction(increment,
		func(input number) int { return input.Value },
		func(input number, value int) number {
			input.Value = value
			return input
		})
}
