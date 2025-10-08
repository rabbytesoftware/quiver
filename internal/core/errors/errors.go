package errors

import "fmt"

type Error struct {
	Code    ErrorCode              `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details"`
}

func Throw(
	code ErrorCode,
	message string,
	details map[string]interface{},
) Error {
	return Error{
		Code:    code,
		Message: message,
		Details: details,
	}
}

func (e Error) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

func (e Error) ShouldRetry() bool {
	return e.Code >= 500
}
