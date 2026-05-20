package chain

import "context"

// BranchAction controls workflow branching after its Run method succeeds.
// It extends Action with custom directions and a direction selector.
type BranchAction[T any] interface {
	// Name returns the name of the BranchAction.
	Name() string

	// Run executes the branch action.
	Run(ctx context.Context, input T) (output T, err error)

	// Directions returns custom directions that NextDirection can return.
	// Built-in directions are already available in Workflow run plans.
	Directions() []string

	// NextDirection selects the next execution path from the Run output.
	// It is called only if Run succeeds (err == nil).
	NextDirection(ctx context.Context, output T) (direction string, err error)
}

// AsBranchAction extends an existing Action with branching logic.
// The wrapped action keeps its original Name and Run behavior, while branchFunc
// determines the next direction from the wrapped action's output.
//
// The optional directions argument describes custom directions returned by branchFunc.
// Built-in directions such as Success, Failure, and Abort are already available
// in Workflow plans.
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
