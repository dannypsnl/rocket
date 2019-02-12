package rocket

import (
	"fmt"
	"net/http"
)

type Guard interface {
	VerifyRequest() error
}

type VerifyError struct {
	status  int
	message string
}

func AuthError(format string, a ...interface{}) *VerifyError {
	return &VerifyError{
		status:  http.StatusForbidden,
		message: fmt.Sprintf(format, a...),
	}
}
func ValidateError(format string, a ...interface{}) *VerifyError {
	return &VerifyError{
		status:  http.StatusBadRequest,
		message: fmt.Sprintf(format, a...),
	}
}

func (v *VerifyError) Status() int {
	return v.status
}
func (v *VerifyError) Error() string {
	return v.message
}
