package response_test

import (
	"net/http"
	"testing"

	"github.com/dannypsnl/rocket/cookie"
	"github.com/dannypsnl/rocket/response"

	asserter "github.com/dannypsnl/assert"
)

type fakeHeaderCounter struct {
	count int
}

func (w *fakeHeaderCounter) Header() http.Header {
	w.count++
	return http.Header(map[string][]string{})
}
func (w *fakeHeaderCounter) Write([]byte) (int, error)  { return 1, nil }
func (w *fakeHeaderCounter) WriteHeader(statusCode int) {}

func TestResponse(t *testing.T) {
	assert := asserter.NewTester(t)

	res := response.New("").
		Headers(response.Headers{
			"x-testing": "hello",
		}).
		Cookies(
			cookie.New("testing", "hello"),
		)
	w := &fakeHeaderCounter{count: 0}
	res.WriteTo(w)
	// include content-type, set header, set cookie
	assert.Eq(w.count, 3)
}

type fakeWriteCounter struct {
	count int
}

func (w *fakeWriteCounter) Header() http.Header {
	return http.Header(map[string][]string{})
}
func (w *fakeWriteCounter) Write([]byte) (int, error) {
	w.count++
	return 1, nil
}
func (w *fakeWriteCounter) WriteHeader(statusCode int) {}

func TestHTTPPipelining(t *testing.T) {
	assert := asserter.NewTester(t)

	testCases := []struct {
		name               string
		expectedWriteTimes int
		keepFunc           func(w http.ResponseWriter)
	}{
		{
			name:               "no pipelining at least would write once",
			expectedWriteTimes: 1,
			keepFunc:           nil,
		},
		{
			name:               "pipelining would keeping write data after response body flush",
			expectedWriteTimes: 3,
			keepFunc: func(w http.ResponseWriter) {
				w.Write([]byte{})
				w.Write([]byte{})
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			w := &fakeWriteCounter{count: 0}
			res := response.New("")
			if testCase.keepFunc != nil {
				res.Keep(testCase.keepFunc)
			}
			res.WriteTo(w)
			assert.Eq(w.count, testCase.expectedWriteTimes)
		})
	}
}

func TestStatusCodeCheck(t *testing.T) {
	t.Run("ChangeStatusCodeTwice", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("should panic when status code be rewritten")
			}
		}()
		response.New("").
			Status(http.StatusNotFound).
			Status(http.StatusOK)
	})
	t.Run("InvalidStatusCode", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("should panic when input code is invalid status code")
			}
		}()
		response.New("").
			Status(1) // NOTE: status code should be a three-digit integer
	})
}
