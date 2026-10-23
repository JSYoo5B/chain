package parameters

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/JSYoo5B/chain"
)

// This assignment checks the generated public signature from a separate file.
var _ func(int, int, func(int) error, ...Option) *chain.Workflow[int] = Build

func TestConstructorArguments(t *testing.T) {
	callbackFailure := errors.New("callback failed")
	cases := []struct {
		name   string
		build  func(func(int) error) *chain.Workflow[int]
		offset int
		err    error
	}{
		{"no options", func(callback func(int) error) *chain.Workflow[int] {
			return Build(2, 3, callback)
		}, 5, nil},
		{"explicit options and callback error", func(callback func(int) error) *chain.Workflow[int] {
			return Build(2, 3, callback, add(3), multiply(2))
		}, 16, callbackFailure},
		{"expanded options preserve order", func(callback func(int) error) *chain.Workflow[int] {
			options := []Option{multiply(2), add(3)}
			return Build(2, 3, callback, options...)
		}, 13, nil},
		{"nil options", func(callback func(int) error) *chain.Workflow[int] {
			var options []Option
			return Build(2, 3, callback, options...)
		}, 5, nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var callbackInputs []int
			w := tc.build(func(output int) error {
				callbackInputs = append(callbackInputs, output)
				return tc.err
			})
			if len(callbackInputs) != 0 {
				t.Fatal("constructor invoked the action callback before Run")
			}
			if err := w.ValidateGraph(); err != nil {
				t.Fatal(err)
			}
			for _, input := range []int{10, 20} {
				output, err := w.Run(context.Background(), input)
				if output != input+tc.offset || !errors.Is(err, tc.err) {
					t.Fatalf("Run(%d) = (%d, %v), want (%d, %v)", input, output, err, input+tc.offset, tc.err)
				}
			}
			want := []int{10 + tc.offset, 20 + tc.offset}
			if !reflect.DeepEqual(callbackInputs, want) {
				t.Fatalf("callback inputs = %v, want %v", callbackInputs, want)
			}
		})
	}
}
