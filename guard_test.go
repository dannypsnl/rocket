package rocket_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dannypsnl/rocket"

	asserter "github.com/dannypsnl/assert"

	"github.com/gavv/httpexpect"
)

type headerGuard struct {
	Auth *string `header:"Auth"`
}

func (h *headerGuard) VerifyRequest() error {
	if h.Auth != nil && *h.Auth == "user1" {
		return nil
	}
	return rocket.AuthError("not allowed")
}

func TestGuard(t *testing.T) {
	testCases := []struct {
		name           string
		guard          rocket.Guard
		testFunc       func(*httpexpect.Request)
		expectedStatus int
	}{
		{
			name:  "valid request would pass guard",
			guard: &headerGuard{},
			testFunc: func(r *httpexpect.Request) {
				r.WithHeader("Auth", "user1").
					Expect().
					Status(http.StatusOK)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid request won't pass guard",
			guard:          &headerGuard{},
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			rk := rocket.Ignite("").
				Mount("/", rocket.Get("/", func() string { return "" }).
					Guard(testCase.guard))
			ts := httptest.NewServer(rk)
			defer ts.Close()
			e := httpexpect.New(t, ts.URL)

			request := e.GET("/")
			if testCase.testFunc != nil {
				testCase.testFunc(request)
			} else {
				request.Expect().
					Status(testCase.expectedStatus)
			}
		})
	}
}

func TestVerifyError(t *testing.T) {
	assert := asserter.NewTester(t)

	err := rocket.AuthError("auth failed")
	assert.Eq(err.Status(), http.StatusForbidden)
	err = rocket.ValidateError("validate failed")
	assert.Eq(err.Status(), http.StatusBadRequest)
}
