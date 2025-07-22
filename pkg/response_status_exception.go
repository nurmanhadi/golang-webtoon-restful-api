package pkg

import "fmt"

type ErrorResponse struct {
	Code    int
	Message string
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("code: %d, error: %s", e.Code, e.Message)
}
func ResponseStatusException(code int, message string) error {
	return &ErrorResponse{
		Code:    code,
		Message: message,
	}
}
