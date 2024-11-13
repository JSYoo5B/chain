package railway

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCompositePipeline(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)

	t.Run("simple composite pipeline", func(t *testing.T) {
		compositePipeline := NewPipeline(
			"CompositeTest",
			NewNumAction(&IncAction{"action1"}),
			NewMsgAction(&GreetAction{"action2"}),
			NewNumAction(&IncAction{"action3"}),
			NewNumAction(&IncAction{"action4"}),
			NewMsgAction(&GreetAction{"action5"}),
		)

		input := TestType{number: 10, message: "f"}
		// {10, f} -> {11, f} -> {11, fo} -> {12, fo} -> {13, fo} -> {13, foo}
		output, direction, err := compositePipeline.Run(context.Background(), input)

		assert.NoError(t, err)
		assert.Equal(t, Success, direction)
		assert.Equal(t, 13, output.number)
		assert.Equal(t, "foo", output.message)
	})

	t.Run("complex composite pipeline", func(t *testing.T) {
		inc2Action := NewPipeline("Inc2Action", &IncAction{"inc1"}, &IncAction{"inc2"})
		greetAction := NewPipeline("GreetAction", &GreetAction{}, &GreetAction{})
		inc5Action := NewPipeline("Inc5Action", &IncAction{"inc3"}, &IncAction{"inc4"}, &IncAction{"inc5"}, &IncAction{"inc"}, &IncAction{"inc7"})

		compositePipeline := NewPipeline("CompositeTest",
			NewNumAction(inc2Action),
			NewMsgAction(greetAction),
			NewNumAction(inc5Action),
		)

		input := TestType{number: 10, message: "f"}
		// {10, f} -> {11, f} -> {12, f}
		// -> {12, fo} -> {12, foo}
		// -> {13, foo} -> {14, foo} -> {15, foo} -> {16, foo} -> {17, foo}
		output, direction, err := compositePipeline.Run(context.Background(), input)

		assert.NoError(t, err)
		assert.Equal(t, Success, direction)
		assert.Equal(t, 17, output.number)
		assert.Equal(t, "foo", output.message)
	})
}

type TestType struct {
	number  int
	message string
}

type NumAction struct {
	action Action[int]
}

func NewNumAction(action Action[int]) *NumAction { return &NumAction{action: action} }
func (n NumAction) Name() string                 { return n.action.Name() }
func (n NumAction) Directions() []string         { return n.action.Directions() }
func (n NumAction) Run(ctx context.Context, input TestType) (TestType, string, error) {
	output := input

	actualInput := input.number
	actualOutput, direction, err := n.action.Run(ctx, actualInput)
	output.number = actualOutput

	return output, direction, err
}

type MsgAction struct {
	action Action[string]
}

func NewMsgAction(action Action[string]) *MsgAction { return &MsgAction{action: action} }
func (m MsgAction) Name() string                    { return m.action.Name() }
func (m MsgAction) Directions() []string            { return m.action.Directions() }
func (m MsgAction) Run(ctx context.Context, input TestType) (TestType, string, error) {
	output := input

	actualInput := input.message
	actualOutput, direction, err := m.action.Run(ctx, actualInput)
	output.message = actualOutput

	return output, direction, err
}

type IncAction struct{ name string }

func (i IncAction) Name() string       { return i.name }
func (IncAction) Directions() []string { return []string{Success, Error, Abort} }
func (IncAction) Run(_ context.Context, input int) (output int, direction string, err error) {
	return input + 1, Success, nil
}

type GreetAction struct{ name string }

func (g GreetAction) Name() string       { return g.name }
func (GreetAction) Directions() []string { return []string{Success, Error, Abort} }
func (GreetAction) Run(_ context.Context, input string) (output string, direction string, err error) {
	return input + "o", Success, nil
}
