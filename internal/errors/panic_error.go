package errors

import (
	"errors"
	"fmt"
)

type PanicError struct {
	Value  any
	Runner string
	Err    error
}

func (e *PanicError) Error() string {
	return fmt.Sprintf("panic in `%s`: %v", e.Runner, e.Value)
}

func (e *PanicError) Unwrap() error {
	return e.Err
}

func NewPanicError(runner string, value any) *PanicError {
	var err error
	switch x := value.(type) {
	case string:
		err = errors.New(x)
	case error:
		err = x
	default:
		err = fmt.Errorf("unknown panic type: %T", value)
	}

	return &PanicError{
		Value:  value,
		Runner: runner,
		Err:    err,
	}
}
