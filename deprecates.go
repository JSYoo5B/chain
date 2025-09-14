package chain

// Deprecated. Use InternalTypeGetter
type AggregateGetter[T any, U any] InternalTypeGetter[T, U]

// Deprecated. Use ExternalTypeSetter
type AggregateSetter[T any, U any] ExternalTypeSetter[T, U]

// Deprecated. Use AdaptAction
func NewAggregateAction[T any, U any](
	action Action[U],
	getter AggregateGetter[T, U],
	setter AggregateSetter[T, U],
) Action[T] {
	return AdaptAction[T, U](
		action,
		(InternalTypeGetter[T, U])(getter),
		(ExternalTypeSetter[T, U])(setter),
	)
}

// Deprecated. Use AsParallelMapAction
func NewParallelMapPipeline[K comparable, T any](name string, action Action[T]) Action[map[K]T] {
	return AsParallelMapAction[K, T](name, action)
}

// Deprecated. Use AsParallelSliceAction
func NewParallelSlicePipeline[T any](name string, action Action[T]) Action[[]T] {
	return AsParallelSliceAction[T](name, action)
}

// Deprecated. Use AsSequenceMapAction
func NewSequenceMapPipeline[K comparable, T any](name string, action Action[T]) Action[map[K]T] {
	return AsSequenceMapAction[K, T](name, action)
}

// Deprecated. Use AsSequenceSliceAction
func NewSequenceSlicePipeline[T any](name string, action Action[T], stopOnError bool) Action[[]T] {
	return AsSequenceSliceAction[T](name, action, stopOnError)
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

// Deprecated. Use AsParallelMapAction
func NewParallelMapAction[K comparable, T any](name string, action Action[T]) Action[map[K]T] {
	return AsParallelMapAction[K, T](name, action)
}

// Deprecated. Use AsParallelSliceAction
func NewParallelSliceAction[T any](name string, action Action[T]) Action[[]T] {
	return AsParallelSliceAction[T](name, action)
}

// Deprecated. Use AsSequenceMapAction
func NewSequenceMapAction[K comparable, T any](name string, action Action[T]) Action[map[K]T] {
	return AsSequenceMapAction[K, T](name, action)
}

// Deprecated. Use AsSequenceSliceAction
func NewSequenceSliceAction[T any](name string, action Action[T], stopOnError bool) Action[[]T] {
	return AsSequenceSliceAction[T](name, action, stopOnError)
}

// Deprecated. Use InternalTypeGetter
type InternalValueGetter[T any, U any] InternalTypeGetter[T, U]

// Deprecated. Use ExternalTypeSetter
type InternalValueSetter[T any, U any] ExternalTypeSetter[T, U]

// Deprecated. Use AdaptAction
func NewTypeAdapterAction[T any, U any](
	action Action[U],
	getter InternalValueGetter[T, U],
	setter InternalValueSetter[T, U],
) Action[T] {
	return AdaptAction[T, U](
		action,
		(InternalTypeGetter[T, U])(getter),
		(ExternalTypeSetter[T, U])(setter),
	)
}

// Deprecated. Use Failure
const Error = Failure
