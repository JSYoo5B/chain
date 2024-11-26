package railway

import "context"

// AggregateGetter extracts a subpart (U) from a composite data structure (T).
// It allows access to a part of the structure without modifying the entire one.
type AggregateGetter[T any, U any] func(T) U

// AggregateSetter updates a composite data structure (T) with a new subpart (U).
// It reintegrates the modified part back into the structure and returns the updated structure.
type AggregateSetter[T any, U any] func(T, U) T

// NewAggregateAction creates an Action that works with a composite data structure (T),
// where T is a complex type (e.g., a struct with multiple fields) and U is the
// data type that the Action operates on. The AggregateGetter and AggregateSetter
// functions are used to extract U from T and re-integrate the processed result back into T.
func NewAggregateAction[T any, U any](
	action Action[U],
	getter AggregateGetter[T, U],
	setter AggregateSetter[T, U],
) Action[T] {
	return &aggregateAction[T, U]{
		action: action,
		getter: getter,
		setter: setter,
	}
}

type aggregateAction[T any, U any] struct {
	action Action[U]
	getter AggregateGetter[T, U]
	setter AggregateSetter[T, U]
}

func (a aggregateAction[T, U]) Name() string { return a.action.Name() }
func (a aggregateAction[T, U]) Run(ctx context.Context, input T) (output T, err error) {
	output = input

	actualInput := a.getter(input)
	actualOutput, err := a.action.Run(ctx, actualInput)
	output = a.setter(output, actualOutput)

	return output, err
}
