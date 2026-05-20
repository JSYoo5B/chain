package chain

import "context"

// BranchAction is an interface for actions that control branching in the execution flow
// of a Workflow. It extends the Action interface and adds methods for handling conditional
// branching based on the execution results.
type BranchAction[T any] interface {
	// Name returns the name of the BranchAction.
	Name() string

	// Run executes the branch action, optionally modifying the input and returning an output.
	// If the input doesn't need changes, it can be passed through as output. The method also
	// returns an error if the action cannot be executed successfully.
	Run(ctx context.Context, input T) (output T, err error)

	// Directions return a list of possible directions that the Workflow can take.
	// These directions are used for validation and must include all possible values that
	// NextDirection can return.
	Directions() []string

	// NextDirection determines the next execution path based on the result of Run.
	// It is called only if Run succeeds (err == nil).
	// The method returns a direction from the list defined by Directions.
	NextDirection(ctx context.Context, output T) (direction string, err error)
}

// AsBranchAction extends an existing Action with branching logic.
// The wrapped action keeps its original Name and Run behavior, while branchFunc determines
// the next direction from the wrapped action's output.
//
// The optional directions argument describes custom directions returned by branchFunc.
// Built-in directions such as Success, Failure, and Abort are already available in Workflow
// plans, so directions may be omitted when branchFunc only returns built-in directions.
func AsBranchAction[T any](baseAction Action[T], branchFunc BranchFunc[T], directions ...string) BranchAction[T] {
	copiedDirections := make([]string, len(directions))
	copy(copiedDirections, directions)

	return &branchActionAdapter[T]{
		Action:     baseAction,
		directions: copiedDirections,
		branchFunc: branchFunc,
	}
}

type branchActionAdapter[T any] struct {
	Action[T]
	directions []string
	branchFunc BranchFunc[T]
}

func (b branchActionAdapter[T]) Directions() []string {
	directions := make([]string, len(b.directions))
	copy(directions, b.directions)
	return directions
}
func (b branchActionAdapter[T]) NextDirection(ctx context.Context, output T) (string, error) {
	return b.branchFunc(ctx, output)
}
