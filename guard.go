package rocket

import (
	"net/http"
)

type Action uint8

const (
	_ Action = iota
	Success
	Failure
	Forward
)

type Guard interface {
	VerifyRequest(*http.Request) (Action, error)
}
