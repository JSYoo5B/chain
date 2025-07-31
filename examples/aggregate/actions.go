package aggregate

import (
	"context"
	"github.com/JSYoo5B/chain"
	"github.com/sirupsen/logrus"
)

// The following functions create two types of actions, Action[int] and Action[string],
// where the generic types are different, making it impossible to handle them within a single pipeline.

func newIncAction(name string) chain.Action[int] {
	runFunc := func(ctx context.Context, input int) (output int, err error) {
		logrus.WithContext(ctx).Infof("Increasing %d to %d", input, input+1)
		return input + 1, nil
	}
	return chain.NewSimpleAction(name, runFunc)
}

func newAppendAction(name string) chain.Action[string] {
	runFunc := func(ctx context.Context, input string) (output string, err error) {
		logrus.WithContext(ctx).Infof("Appending %s with trailing o", input)
		return input + "o", nil
	}
	return chain.NewSimpleAction(name, runFunc)
}

// By aggregating int and string into a single `Pair` struct,
// new actions are defined for `Action[Pair]` which combine each `Action[int]` and `Action[string]`.
// This enables handling different types of actions within a single pipeline.

type Pair struct {
	number  int
	message string
}

func numberToPair(action chain.Action[int]) chain.Action[Pair] {
	getter := func(c Pair) int { return c.number }
	setter := func(c Pair, i int) Pair {
		c.number = i
		return c
	}
	return chain.NewAggregateAction(action, getter, setter)
}

func messageToPair(action chain.Action[string]) chain.Action[Pair] {
	getter := func(c Pair) string { return c.message }
	setter := func(c Pair, s string) Pair {
		c.message = s
		return c
	}
	return chain.NewAggregateAction(action, getter, setter)
}
