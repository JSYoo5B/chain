package branch

import (
	"context"

	"github.com/JSYoo5B/chain"
)

func checkSign() chain.BranchAction[int] {
	return chain.NewSimpleBranchAction[int]("sign", nil, []string{"negative", "nonnegative"},
		func(_ context.Context, output int) (string, error) {
			if output < 0 {
				return "negative", nil
			}
			return "nonnegative", nil
		})
}

func incrementAndCheckParity() chain.BranchAction[int] {
	return chain.AsBranchAction(inc(), func(_ context.Context, output int) (string, error) {
		if output%2 == 0 {
			return "even", nil
		}
		return "odd", nil
	}, "even", "odd")
}
