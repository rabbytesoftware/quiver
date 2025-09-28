package errors

import (
	"fmt"
)

type Err struct {
	Code    int
	Message string
}

func NewErr(code int, message string) *Err {
	return &Err{
		Code:    code,
		Message: message,
	}
}

func (e *Err) Error() error {
	return fmt.Errorf("%s", e.Message)
}
