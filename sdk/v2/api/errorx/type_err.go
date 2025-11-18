package errorx

import "fmt"

type TypeMismatchCause struct {
	ActualType   string
	ExpectedType string
	ActualValue  interface{}
}

func (e *TypeMismatchCause) Error() string {
	return fmt.Sprintf("type mismatch: expected=%s actual=%s value=%v",
		e.ExpectedType, e.ActualType, e.ActualValue)
}
