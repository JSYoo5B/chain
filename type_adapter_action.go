package chain

import "context"

// InternalValueGetter extracts a subpart (U) from a composite data structure (T).
// It allows access to a part of the structure without modifying the entire one.
type InternalValueGetter[T any, U any] func(T) U

// InternalValueSetter updates a composite data structure (T) with a new subpart (U).
// It reintegrates the modified part back into the structure and returns the updated structure.
type InternalValueSetter[T any, U any] func(T, U) T

// NewTypeAdapterAction creates an Action that works with a composite data structure (T),
// where T is a complex type (e.g., a struct with multiple fields) and U is the
// data type that the Action operates on. The InternalValueGetter and InternalValueSetter
// functions are used to extract U from T and re-integrate the processed result back into T.
func NewTypeAdapterAction[T any, U any](
	action Action[U],
	getter InternalValueGetter[T, U],
	setter InternalValueSetter[T, U],
) Action[T] {
	if action == nil {
		panic("action cannot be nil")
	} else if getter == nil {
		panic("getter cannot be nil")
	} else if setter == nil {
		panic("setter cannot be nil")
	}

	return &typeAdapterAction[T, U]{
		action: action,
		getter: getter,
		setter: setter,
	}
}

type typeAdapterAction[T any, U any] struct {
	action Action[U]
	getter InternalValueGetter[T, U]
	setter InternalValueSetter[T, U]
}

func (a typeAdapterAction[T, U]) Name() string { return a.action.Name() }
func (a typeAdapterAction[T, U]) Run(ctx context.Context, input T) (output T, err error) {
	output = input

	actualInput := a.getter(input)
	actualOutput, err := a.action.Run(ctx, actualInput)
	output = a.setter(output, actualOutput)

	return output, err
}
