package collection

import (
	"context"
	"errors"
	"strings"

	"github.com/JSYoo5B/chain"
)

func double() chain.Action[int] {
	return chain.NewSimpleAction("double", func(_ context.Context, input int) (int, error) {
		return input * 2, nil
	})
}

func sequenceNumbers() chain.Action[[]int] {
	return chain.AsSequenceSliceAction("double each", double(), true)
}

func sequenceWithError(stopOnError bool) chain.Action[[]int] {
	action := chain.NewSimpleAction("double except two", func(_ context.Context, input int) (int, error) {
		if input == 2 {
			return input, errors.New("two is unchanged")
		}
		return input * 2, nil
	})
	return chain.AsSequenceSliceAction("double each", action, stopOnError)
}

func sequenceWords() chain.Action[map[string]string] {
	upper := chain.NewSimpleAction("uppercase", func(_ context.Context, input string) (string, error) {
		return strings.ToUpper(input), nil
	})
	return chain.AsSequenceMapAction[string]("uppercase each", upper)
}

func parallelNumbers() chain.Action[[]int] {
	return chain.AsParallelSliceAction("double each", double())
}

func parallelNumberMap() chain.Action[map[string]int] {
	return chain.AsParallelMapAction[string]("double each", double())
}
