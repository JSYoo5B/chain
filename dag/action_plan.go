package dag

type ActionPlan[T any] map[string]Action[T]

const (
	Success = "success"
	Error   = "error"
	Abort   = "abort"
)

func TerminationPlan[T any]() ActionPlan[T] {
	return nil
}
