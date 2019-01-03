package rocket_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dannypsnl/rocket"

	"github.com/gavv/httpexpect"
)

type headerGuard struct{}

func (h *headerGuard) VerifyRequest(r *http.Request) (rocket.Action, error) {
	if r.Header.Get("Auth") == "user1" {
		return rocket.Success, nil
	}
	return rocket.Failure, errors.New("not allowed")
}

func TestGuard(t *testing.T) {
	rk := rocket.Ignite(":8081").
		Mount("/", rocket.Get("/", func(h *headerGuard) string {
			return "pass"
		}))
	ts := httptest.NewServer(rk)
	defer ts.Close()

	e := httpexpect.New(t, ts.URL)

	e.GET("/").
		Expect().
		Status(http.StatusBadRequest)
	e.GET("/").WithHeader("Auth", "user1").
		Expect().
		Status(http.StatusOK)
}
