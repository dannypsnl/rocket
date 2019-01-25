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

type forwardGuard struct{}

func (f *forwardGuard) VerifyRequest(r *http.Request) (rocket.Action, error) {
	return rocket.Forward, nil
}

func TestGuard(t *testing.T) {
	testCases := []struct {
		name           string
		handlerFunc    interface{}
		testFunc       func(*httpexpect.Request)
		expectedStatus int
	}{
		{
			name:        "valid request would pass guard",
			handlerFunc: func(*headerGuard) string { return "" },
			testFunc: func(r *httpexpect.Request) {
				r.WithHeader("Auth", "user1").
					Expect().
					Status(http.StatusOK)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid request won't pass guard",
			handlerFunc:    func(*headerGuard) string { return "" },
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "action forward is give up decision right",
			handlerFunc:    func(*forwardGuard) string { return "" },
			expectedStatus: http.StatusOK,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			rk := rocket.Ignite("").
				Mount("/", rocket.Get("/", testCase.handlerFunc))
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
