package internal

import (
	"context"
	"github.com/JSYoo5B/dago/railway"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

// The following functions create two types of actions, Action[int] and Action[string],
// where the generic types are different, making it impossible to handle them within a single pipeline.

func newIncAction(name string) railway.Action[int] {
	runFunc := func(_ context.Context, input int) (output int, err error) {
		return input + 1, nil
	}
	return railway.NewSimpleAction(name, runFunc)
}

func newAppendAction(name string) railway.Action[string] {
	runFunc := func(_ context.Context, input string) (output string, err error) {
		return input + "o", nil
	}
	return railway.NewSimpleAction(name, runFunc)
}

// By aggregating int and string into a single `Wrap` struct,
// new actions are defined for `Action[Wrap]` which combine each `Action[int]` and `Action[string]`.
// This enables handling different types of actions within a single pipeline.

type Wrap struct {
	number  int
	message string
}

func newNumberAction(action railway.Action[int]) railway.Action[Wrap] {
	getter := func(c Wrap) int { return c.number }
	setter := func(c Wrap, i int) Wrap {
		c.number = i
		return c
	}
	return railway.NewAggregateAction(action, getter, setter)
}

func newMessageAction(action railway.Action[string]) railway.Action[Wrap] {
	getter := func(c Wrap) string { return c.message }
	setter := func(c Wrap, s string) Wrap {
		c.message = s
		return c
	}
	return railway.NewAggregateAction(action, getter, setter)
}

// This test demonstrates how different types of actions can be handled within a single pipeline.
// It shows that both simple Action and Pipeline (treated as Action) can be combined and processed together.

func TestAggregatePipeline(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)

	t.Run("simple aggregate pipeline", func(t *testing.T) {
		aggregatePipeline := railway.NewPipeline(
			"ActionAggregateTest",
			newNumberAction(newIncAction("action1")),
			newMessageAction(newAppendAction("action2")),
			newNumberAction(newIncAction("action3")),
			newNumberAction(newIncAction("action4")),
			newMessageAction(newAppendAction("action5")),
		)

		input := Wrap{number: 10, message: "f"}
		// {10, f} -> {11, f} -> {11, fo} -> {12, fo} -> {13, fo} -> {13, foo}
		output, err := aggregatePipeline.Run(context.Background(), input)

		assert.NoError(t, err)
		assert.Equal(t, 13, output.number)
		assert.Equal(t, "foo", output.message)
	})

	t.Run("nested pipeline aggregate pipeline", func(t *testing.T) {
		inc2Action := railway.NewPipeline(
			"Inc2Action",
			newIncAction("inc1"),
			newIncAction("inc2"),
		)
		append2Action := railway.NewPipeline(
			"Append2Action",
			newAppendAction("append1"),
			newAppendAction("append2"),
		)
		inc5Action := railway.NewPipeline(
			"Inc5Action",
			newIncAction("inc3"),
			newIncAction("inc4"),
			newIncAction("inc5"),
			newIncAction("inc6"),
			newIncAction("inc7"),
		)

		aggregatePipeline := railway.NewPipeline(
			"PipelineAggregateTest",
			newNumberAction(inc2Action),
			newMessageAction(append2Action),
			newNumberAction(inc5Action),
		)

		input := Wrap{number: 10, message: "f"}
		// {10, f} -> {11, f} -> {12, f}
		// -> {12, fo} -> {12, foo}
		// -> {13, foo} -> {14, foo} -> {15, foo} -> {16, foo} -> {17, foo}
		output, err := aggregatePipeline.Run(context.Background(), input)

		assert.NoError(t, err)
		assert.Equal(t, 17, output.number)
		assert.Equal(t, "foo", output.message)
	})
}
