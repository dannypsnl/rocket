package rocket

import (
	"fmt"
	"net/http"

	"github.com/dannypsnl/rocket/response"
)

// Guard is an interface that context can implement, when context implement this, context can reject request with a *response.Response.
type Guard interface {
	VerifyRequest() *response.Response
}

func AuthError(format string, a ...interface{}) *response.Response {
	return response.New(fmt.Sprintf(format, a...)).
		Status(http.StatusForbidden)
}
func ValidateError(format string, a ...interface{}) *response.Response {
	return response.New(fmt.Sprintf(format, a...)).
		Status(http.StatusBadRequest)
}
