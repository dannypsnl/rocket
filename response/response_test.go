package response_test

import (
	"net/http"
	"testing"

	"github.com/dannypsnl/rocket/cookie"
	"github.com/dannypsnl/rocket/response"

	asserter "github.com/dannypsnl/assert"
)

type fakeRespWriter struct {
	count int
}

func (w *fakeRespWriter) Header() http.Header {
	w.count++
	return http.Header(map[string][]string{})
}
func (w *fakeRespWriter) Write([]byte) (int, error)  { return 1, nil }
func (w *fakeRespWriter) WriteHeader(statusCode int) {}

func TestResponse(t *testing.T) {
	assert := asserter.NewTester(t)

	res := response.New("").
		Headers(response.Headers{
			"x-testing": "hello",
		}).
		Cookies(
			cookie.New("testing", "hello"),
		)
	w := &fakeRespWriter{count: 0}
	res.WriteTo(w)
	// include content-type, set header, set cookie
	assert.Eq(w.count, 3)
}

type fakeWriteCounter struct {
	count int
}

func (w *fakeWriteCounter) Flush() {}
func (w *fakeWriteCounter) CloseNotify() <-chan bool {
	return make(<-chan bool)
}
func (w *fakeWriteCounter) Header() http.Header {
	return http.Header(map[string][]string{})
}
func (w *fakeWriteCounter) Write([]byte) (int, error) {
	w.count++
	return 1, nil
}
func (w *fakeWriteCounter) WriteHeader(statusCode int) {}

func TestHTTPStreaming(t *testing.T) {
	assert := asserter.NewTester(t)

	testCases := []struct {
		name               string
		expectedWriteTimes int
		streamFunc         response.KeepFunc
	}{
		{
			name:               "nil func should be ignore",
			expectedWriteTimes: 1,
			streamFunc:         nil,
		},
		{
			name:               "streaming would keeping write data after response body flush",
			expectedWriteTimes: 3,
			streamFunc: func(w http.ResponseWriter) bool {
				w.Write([]byte{})
				w.Write([]byte{})
				return false
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			w := &fakeWriteCounter{count: 0}
			res := response.Stream(testCase.streamFunc)
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
