package chain

// RunPlan represents a map that associates a direction
// (Success, Failure, Abort, and other custom branching directions) with the next Action to execute.
// It is used to define the flow of actions in a Workflow based on the direction of execution.
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

// TerminationPlan returns a RunPlan with all directions leading to termination immediately,
// providing a clear sign of termination rather than using nil.
func TerminationPlan[T any]() RunPlan[T] {
	return nil
}

// SuccessOnlyPlan returns a RunPlan where only a success direction has a valid next action,
// and Failure and Abort both lead to termination.
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

// DefaultPlanWithAbort returns a RunPlan with valid next actions for Success, Failure, and Abort.
func DefaultPlanWithAbort[T any](success, error, abort Action[T]) RunPlan[T] {
	return RunPlan[T]{
		Success: success,
		Failure: error,
		Abort:   abort,
	}
}

// Terminate provides an Action that explicitly stops the execution of a Workflow.
//
// When returned from a RunPlan, it signals that the Workflow should halt
// and no further actions should be-executed. This serves as a clear, intentional
// way to stop a Workflow, as opposed to returning raw nil.
func Terminate[T any]() Action[T] {
	return nil
}
