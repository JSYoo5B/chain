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

// deprecated. Use NewParallelMapAction
func NewParallelMapPipeline[K comparable, T any](name string, action Action[T]) Action[map[K]T] {
	return NewParallelMapAction[K, T](name, action)
}

// deprecated. Use NewParallelSliceAction
func NewParallelSlicePipeline[T any](name string, action Action[T]) Action[[]T] {
	return NewParallelSliceAction[T](name, action)
}

// deprecated. Use NewSequenceMapAction
func NewSequenceMapPipeline[K comparable, T any](name string, action Action[T]) Action[map[K]T] {
	return NewSequenceMapAction[K, T](name, action)
}

// deprecated. Use NewSequenceSliceAction
func NewSequenceSlicePipeline[T any](name string, action Action[T], stopOnError bool) Action[[]T] {
	return NewSequenceSliceAction[T](name, action, stopOnError)
}
