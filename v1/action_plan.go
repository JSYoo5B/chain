package v1

type ActionPlan[T any] map[string]Action[T]

const (
	Success = "success"
	Error   = "error"
	Abort   = "abort"
)

func TerminationPlan[T any]() ActionPlan[T] {
	return nil
}

func SuccessOnlyPlan[T any](success Action[T]) ActionPlan[T] {
	return ActionPlan[T]{
		Success: success,
		Error:   Terminate[T](),
		Abort:   Terminate[T](),
	}
}

func DefaultPlan[T any](success, error Action[T]) ActionPlan[T] {
	return ActionPlan[T]{
		Success: success,
		Error:   error,
		Abort:   Terminate[T](),
	}
}

func DefaultPlanWithAbort[T any](success, error, abort Action[T]) ActionPlan[T] {
	return ActionPlan[T]{
		Success: success,
		Error:   error,
		Abort:   abort,
	}
}
