package adapter

import (
	"context"
	"strings"

	"github.com/JSYoo5B/chain"
)

func newNormalizeCustomerAction(name string) chain.Action[string] {
	return chain.NewSimpleAction(
		name,
		func(_ context.Context, input string) (string, error) {
			return strings.Join(strings.Fields(input), " "), nil
		})
}

func newAddPriceAction(name string, cents int) chain.Action[int] {
	return chain.NewSimpleAction(
		name,
		func(_ context.Context, input int) (int, error) {
			return input + cents, nil
		})
}

type cart struct {
	customer   string
	totalCents int
}

func customerToCart(action chain.Action[string]) chain.Action[cart] {
	return chain.AdaptAction(
		action,
		func(c cart) string { return c.customer },
		func(c cart, customer string) cart {
			c.customer = customer
			return c
		},
	)
}

func totalCentsToCart(action chain.Action[int]) chain.Action[cart] {
	return chain.AdaptAction(
		action,
		func(c cart) int { return c.totalCents },
		func(c cart, totalCents int) cart {
			c.totalCents = totalCents
			return c
		},
	)
}
