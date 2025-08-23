package chain

// Deprecated. Use InternalValueGetter
type AggregateGetter[T any, U any] InternalValueGetter[T, U]

// Deprecated. Use InternalValueSetter
type AggregateSetter[T any, U any] InternalValueSetter[T, U]

// Deprecated. Use NewTypeAdapterAction
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

// Deprecated. Use NewParallelMapAction
func NewParallelMapPipeline[K comparable, T any](name string, action Action[T]) Action[map[K]T] {
	return NewParallelMapAction[K, T](name, action)
}

// Deprecated. Use NewParallelSliceAction
func NewParallelSlicePipeline[T any](name string, action Action[T]) Action[[]T] {
	return NewParallelSliceAction[T](name, action)
}

// Deprecated. Use NewSequenceMapAction
func NewSequenceMapPipeline[K comparable, T any](name string, action Action[T]) Action[map[K]T] {
	return NewSequenceMapAction[K, T](name, action)
}

// Deprecated. Use NewSequenceSliceAction
func NewSequenceSlicePipeline[T any](name string, action Action[T], stopOnError bool) Action[[]T] {
	return NewSequenceSliceAction[T](name, action, stopOnError)
}

// Deprecated. Use Workflow
type Pipeline[T any] struct {
	*Workflow[T]
}

// Deprecated. Use NewWorkflow
func NewPipeline[T any](name string, memberActions ...Action[T]) *Pipeline[T] {
	return &Pipeline[T]{
		Workflow: NewWorkflow[T](name, memberActions...),
	}
}
