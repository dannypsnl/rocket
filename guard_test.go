package rocket_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dannypsnl/rocket"

	"github.com/gavv/httpexpect"
)

type headerGuard struct {
	Auth *string `header:"Auth"`
}

func (h *headerGuard) VerifyRequest() error {
	if h.Auth != nil && *h.Auth == "user1" {
		return nil
	}
	return errors.New("not allowed")
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
			handlerFunc: func() string { return "" },
			testFunc: func(r *httpexpect.Request) {
				r.WithHeader("Auth", "user1").
					Expect().
					Status(http.StatusOK)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid request won't pass guard",
			handlerFunc:    func() string { return "" },
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			rk := rocket.Ignite("").
				Mount("/", rocket.Get("/", testCase.handlerFunc).
					Guard(&headerGuard{}))
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
