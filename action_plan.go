package chain

// ActionPlan represents a map that associates a direction
// (Success, Failure, Abort, and other custom branching directions) with the next Action to execute.
// It is used to define the flow of actions in a Workflow based on the direction of execution.
type ActionPlan[T any] map[string]Action[T]

const (
	// Success represents the direction indicating that the action completed successfully
	// and the Workflow should continue.
	Success = "success"
	// Failure represents the direction indicating that an error occurred,
	// and the Workflow should handle it accordingly.
	Failure = "failure"
	// Abort represents the direction indicating that
	// the Workflow execution should be aborted immediately.
	// This can occur due to a specific Abort condition or
	// in cases of unexpected errors or panics that cause the Workflow to halt.
	Abort = "abort"
)

// TerminationPlan returns an ActionPlan with all directions leading to termination immediately,
// providing a clear sign of termination rather than using nil.
func TerminationPlan[T any]() ActionPlan[T] {
	return nil
}

// SuccessOnlyPlan returns an ActionPlan where only a success direction has a valid next action,
// and Failure and Abort both lead to termination.
func SuccessOnlyPlan[T any](success Action[T]) ActionPlan[T] {
	return ActionPlan[T]{
		Success: success,
		Failure: Terminate[T](),
		Abort:   Terminate[T](),
	}
}

// DefaultPlan returns a standard ActionPlan with valid next actions for Success and Failure,
// and Termination for Abort.
func DefaultPlan[T any](success, error Action[T]) ActionPlan[T] {
	return ActionPlan[T]{
		Success: success,
		Failure: error,
		Abort:   Terminate[T](),
	}
}

// DefaultPlanWithAbort returns an ActionPlan with valid next actions for Success, Failure, and Abort.
func DefaultPlanWithAbort[T any](success, error, abort Action[T]) ActionPlan[T] {
	return ActionPlan[T]{
		Success: success,
		Failure: error,
		Abort:   abort,
	}
}
