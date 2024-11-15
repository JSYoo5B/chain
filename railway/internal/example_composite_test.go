package internal

import (
	"context"
	"github.com/JSYoo5B/dago/railway"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCompositeTypePipeline(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)

	t.Run("simple composite pipeline", func(t *testing.T) {
		compositePipeline := railway.NewPipeline(
			"CompositeTest",
			NewNumberWrappedAction(&IncAction{"action1"}),
			NewMessageWrappedAction(&GreetAction{"action2"}),
			NewNumberWrappedAction(&IncAction{"action3"}),
			NewNumberWrappedAction(&IncAction{"action4"}),
			NewMessageWrappedAction(&GreetAction{"action5"}),
		)

		input := TestType{number: 10, message: "f"}
		// {10, f} -> {11, f} -> {11, fo} -> {12, fo} -> {13, fo} -> {13, foo}
		output, direction, err := compositePipeline.Run(context.Background(), input)

		assert.NoError(t, err)
		assert.Equal(t, railway.Success, direction)
		assert.Equal(t, 13, output.number)
		assert.Equal(t, "foo", output.message)
	})

	t.Run("complex composite pipeline", func(t *testing.T) {
		inc2Action := railway.NewPipeline("Inc2Action", &IncAction{"inc1"}, &IncAction{"inc2"})
		greetAction := railway.NewPipeline("GreetAction", &GreetAction{}, &GreetAction{})
		inc5Action := railway.NewPipeline("Inc5Action", &IncAction{"inc3"}, &IncAction{"inc4"}, &IncAction{"inc5"}, &IncAction{"inc"}, &IncAction{"inc7"})

		compositePipeline := railway.NewPipeline("CompositeTest",
			NewNumberWrappedAction(inc2Action),
			NewMessageWrappedAction(greetAction),
			NewNumberWrappedAction(inc5Action),
		)

		input := TestType{number: 10, message: "f"}
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

type TestType struct {
	number  int
	message string
}

type NumberWrappedAction struct {
	action railway.Action[int]
}

func NewNumberWrappedAction(action railway.Action[int]) *NumberWrappedAction {
	return &NumberWrappedAction{action: action}
}
func (n NumberWrappedAction) Name() string         { return n.action.Name() }
func (n NumberWrappedAction) Directions() []string { return n.action.Directions() }
func (n NumberWrappedAction) Run(ctx context.Context, input TestType) (TestType, string, error) {
	output := input

	actualInput := input.number
	actualOutput, direction, err := n.action.Run(ctx, actualInput)
	output.number = actualOutput

	return output, direction, err
}

type MessageWrappedAction struct {
	action railway.Action[string]
}

func NewMessageWrappedAction(action railway.Action[string]) *MessageWrappedAction {
	return &MessageWrappedAction{action: action}
}
func (m MessageWrappedAction) Name() string         { return m.action.Name() }
func (m MessageWrappedAction) Directions() []string { return m.action.Directions() }
func (m MessageWrappedAction) Run(ctx context.Context, input TestType) (TestType, string, error) {
	output := input

	actualInput := input.message
	actualOutput, direction, err := m.action.Run(ctx, actualInput)
	output.message = actualOutput

	return output, direction, err
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
