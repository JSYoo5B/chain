package chain

// RunPlan maps directions to the next Action to execute.
// Directions can be built-in values such as Success, Failure, and Abort, or
// custom directions returned by a BranchAction.
type RunPlan[T any] map[string]Action[T]

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

// TerminationPlan returns a RunPlan that terminates immediately.
func TerminationPlan[T any]() RunPlan[T] {
	return nil
}

// SuccessOnlyPlan returns a RunPlan where only Success has a next action.
// Failure and Abort both lead to termination.
func SuccessOnlyPlan[T any](success Action[T]) RunPlan[T] {
	return RunPlan[T]{
		Success: success,
		Failure: Terminate[T](),
		Abort:   Terminate[T](),
	}
}

// DefaultPlan returns a standard RunPlan with valid next actions for Success and Failure,
// and Termination for Abort.
func DefaultPlan[T any](success, error Action[T]) RunPlan[T] {
	return RunPlan[T]{
		Success: success,
		Failure: error,
		Abort:   Terminate[T](),
	}
}

// DefaultPlanWithAbort returns a RunPlan for Success, Failure, and Abort.
func DefaultPlanWithAbort[T any](success, error, abort Action[T]) RunPlan[T] {
	return RunPlan[T]{
		Success: success,
		Failure: error,
		Abort:   abort,
	}
}

// Terminate provides an Action that explicitly stops the execution of a Workflow.
//
// When returned from a RunPlan, it signals that the Workflow should halt and no
// further actions should be executed. This is clearer than returning raw nil.
func Terminate[T any]() Action[T] {
	return nil
}
