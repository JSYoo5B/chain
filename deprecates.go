package chain

// deprecated. Use InternalValueGetter
type AggregateGetter[T any, U any] InternalValueGetter[T, U]

// deprecated. Use InternalValueSetter
type AggregateSetter[T any, U any] InternalValueSetter[T, U]

// deprecated. Use NewTypeAdapterAction
func NewAggregateAction[T any, U any](
	action Action[U],
	getter AggregateGetter[T, U],
	setter AggregateSetter[T, U],
) Action[T] {
	return NewTypeAdapterAction[T, U](
		action,
		(InternalValueGetter[T, U])(getter),
		(InternalValueSetter[T, U])(setter),
	)
}
