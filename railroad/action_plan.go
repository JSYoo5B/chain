package railroad

// ActionPlan represents a map that associates a direction (Success, Error, Abort, and other custom branching directions)
// with the next Action to execute. It is used to define the flow of actions in a pipeline based on the direction of execution.
type ActionPlan[T any] map[string]Action[T]

const (
	// Success represents the direction indicating that the action completed successfully and the pipeline should continue.
	Success = "success"
	// Error represents the direction indicating that an error occurred, and the pipeline should handle it accordingly.
	Error = "error"
	// Abort represents the direction indicating that the pipeline execution should be aborted immediately.
	// This can occur due to a specific Abort condition or in cases of unexpected errors or panics that cause the pipeline to halt.
	Abort = "abort"
)

// TerminationPlan returns an ActionPlan with all directions leading to termination immediately,
// providing a clear indication of termination rather than returning nil.
func TerminationPlan[T any]() ActionPlan[T] {
	return nil
}

// SuccessOnlyPlan returns an ActionPlan where only a success direction has a valid next action,
// and Error and Abort both lead to termination.
func SuccessOnlyPlan[T any](success Action[T]) ActionPlan[T] {
	return ActionPlan[T]{
		Success: success,
		Error:   Terminate[T](),
		Abort:   Terminate[T](),
	}
}

// DefaultPlan returns a standard ActionPlan with valid next actions for Success and Error, and Termination for Abort.
func DefaultPlan[T any](success, error Action[T]) ActionPlan[T] {
	return ActionPlan[T]{
		Success: success,
		Error:   error,
		Abort:   Terminate[T](),
	}
}

// DefaultPlanWithAbort returns an ActionPlan with valid next actions for Success, Error, and Abort.
func DefaultPlanWithAbort[T any](success, error, abort Action[T]) ActionPlan[T] {
	return ActionPlan[T]{
		Success: success,
		Error:   error,
		Abort:   abort,
	}
}
