package rocket_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dannypsnl/rocket"
	"github.com/dannypsnl/rocket/response"

	"github.com/gavv/httpexpect"
)

type headerGuard struct {
	Auth *string `header:"Auth"`
}

func (h *headerGuard) VerifyRequest() *response.Response {
	if h.Auth != nil && *h.Auth == "user1" {
		return nil
	}
	return rocket.AuthError("not allowed")
}

func TestGuard(t *testing.T) {
	testCases := []struct {
		name           string
		testFunc       func(*httpexpect.Request)
		expectedStatus int
	}{
		{
			name: "valid request would pass guard",
			testFunc: func(r *httpexpect.Request) {
				r.WithHeader("Auth", "user1").
					Expect().
					Status(http.StatusOK)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid request won't pass guard",
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			rk := rocket.Ignite(-1).
				Mount(rocket.Get("/", func(*headerGuard) string { return "" }))
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
