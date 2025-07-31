package aggregate

import (
	"context"
	"github.com/JSYoo5B/chain"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAggregatePipeline(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)

	t.Run("simple aggregate pipeline", func(t *testing.T) {
		aggregatePipeline := chain.NewPipeline(
			"ActionAggregateTest",
			numberToPair(newIncAction("action1")),
			messageToPair(newAppendAction("action2")),
			numberToPair(newIncAction("action3")),
			numberToPair(newIncAction("action4")),
			messageToPair(newAppendAction("action5")),
		)

		input := Pair{number: 10, message: "f"}
		// {10, f} -> {11, f} -> {11, fo} -> {12, fo} -> {13, fo} -> {13, foo}
		output, err := aggregatePipeline.Run(context.Background(), input)

		assert.NoError(t, err)
		assert.Equal(t, 13, output.number)
		assert.Equal(t, "foo", output.message)
	})

	t.Run("nested pipeline aggregate pipeline", func(t *testing.T) {
		inc2Action := chain.NewPipeline(
			"Inc2Action",
			newIncAction("inc1"),
			newIncAction("inc2"),
		)
		append2Action := chain.NewPipeline(
			"Append2Action",
			newAppendAction("append1"),
			newAppendAction("append2"),
		)
		inc5Action := chain.NewPipeline(
			"Inc5Action",
			newIncAction("inc3"),
			newIncAction("inc4"),
			newIncAction("inc5"),
			newIncAction("inc6"),
			newIncAction("inc7"),
		)

		aggregatePipeline := chain.NewPipeline(
			"PipelineAggregateTest",
			numberToPair(inc2Action),
			messageToPair(append2Action),
			numberToPair(inc5Action),
		)

		input := Pair{number: 10, message: "f"}
		// {10, f} -> {11, f} -> {12, f}
		// -> {12, fo} -> {12, foo}
		// -> {13, foo} -> {14, foo} -> {15, foo} -> {16, foo} -> {17, foo}
		output, err := aggregatePipeline.Run(context.Background(), input)

		assert.NoError(t, err)
		assert.Equal(t, 17, output.number)
		assert.Equal(t, "foo", output.message)
	})
}
