package internal

import (
	"context"
	"github.com/JSYoo5B/dago/railway"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAggregatePipeline(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)

	t.Run("simple aggregate pipeline", func(t *testing.T) {
		compositePipeline := railway.NewPipeline(
			"CompositeTest",
			NewNumberAction(&IncAction{"action1"}),
			NewMessageAction(&GreetAction{"action2"}),
			NewNumberAction(&IncAction{"action3"}),
			NewNumberAction(&IncAction{"action4"}),
			NewMessageAction(&GreetAction{"action5"}),
		)

		input := Aggregate{number: 10, message: "f"}
		// {10, f} -> {11, f} -> {11, fo} -> {12, fo} -> {13, fo} -> {13, foo}
		output, direction, err := compositePipeline.Run(context.Background(), input)

		assert.NoError(t, err)
		assert.Equal(t, railway.Success, direction)
		assert.Equal(t, 13, output.number)
		assert.Equal(t, "foo", output.message)
	})

	t.Run("complex aggregate pipeline", func(t *testing.T) {
		inc2Action := railway.NewPipeline("Inc2Action", &IncAction{"inc1"}, &IncAction{"inc2"})
		greetAction := railway.NewPipeline("GreetAction", &GreetAction{}, &GreetAction{})
		inc5Action := railway.NewPipeline("Inc5Action", &IncAction{"inc3"}, &IncAction{"inc4"}, &IncAction{"inc5"}, &IncAction{"inc"}, &IncAction{"inc7"})

		compositePipeline := railway.NewPipeline("CompositeTest",
			NewNumberAction(inc2Action),
			NewMessageAction(greetAction),
			NewNumberAction(inc5Action),
		)

		input := Aggregate{number: 10, message: "f"}
		// {10, f} -> {11, f} -> {12, f}
		// -> {12, fo} -> {12, foo}
		// -> {13, foo} -> {14, foo} -> {15, foo} -> {16, foo} -> {17, foo}
		output, direction, err := compositePipeline.Run(context.Background(), input)

		assert.NoError(t, err)
		assert.Equal(t, railway.Success, direction)
		assert.Equal(t, 17, output.number)
		assert.Equal(t, "foo", output.message)
	})
}

type Aggregate struct {
	number  int
	message string
}

func NewNumberAction(action railway.Action[int]) railway.Action[Aggregate] {
	getter := func(c Aggregate) int { return c.number }
	setter := func(c Aggregate, i int) Aggregate {
		c.number = i
		return c
	}
	return railway.NewAggregateAction(action, getter, setter)
}

func NewMessageAction(action railway.Action[string]) railway.Action[Aggregate] {
	getter := func(c Aggregate) string { return c.message }
	setter := func(c Aggregate, s string) Aggregate {
		c.message = s
		return c
	}
	return railway.NewAggregateAction(action, getter, setter)
}

type IncAction struct{ name string }

func (i IncAction) Name() string { return i.name }
func (IncAction) Directions() []string {
	return []string{railway.Success, railway.Error, railway.Abort}
}
func (IncAction) Run(_ context.Context, input int) (output int, direction string, err error) {
	return input + 1, railway.Success, nil
}

type GreetAction struct{ name string }

func (g GreetAction) Name() string { return g.name }
func (GreetAction) Directions() []string {
	return []string{railway.Success, railway.Error, railway.Abort}
}
func (GreetAction) Run(_ context.Context, input string) (output string, direction string, err error) {
	return input + "o", railway.Success, nil
}
