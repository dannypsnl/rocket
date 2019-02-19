package response_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dannypsnl/rocket/cookie"
	"github.com/dannypsnl/rocket/response"
	"github.com/gavv/httpexpect"

	asserter "github.com/dannypsnl/assert"
)

func TestRespAsHTTPHandler(t *testing.T) {
	testCases := []struct {
		name    string
		resp    *response.Response
		expectF func(r *httpexpect.Request)
	}{
		{
			name: "say hello",
			resp: response.New("hello"),
			expectF: func(r *httpexpect.Request) {
				r.Expect().Body().Contains("hello")
			},
		},
		{
			name: "header and cookie",
			resp: response.New("").
				Headers(response.Headers{
					"x-testing": "hello",
				}).
				Cookies(
					cookie.New("testing", "hello"),
				),
			expectF: func(r *httpexpect.Request) {
				r.Expect().Header("x-testing").Equal("hello")
				r.Expect().Cookie("testing").Value().Equal("hello")
			},
		},
		{
			name: "status code",
			resp: response.New("").Status(http.StatusNotFound),
			expectF: func(r *httpexpect.Request) {
				r.Expect().Status(http.StatusNotFound)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ts := httptest.NewServer(testCase.resp)
			defer ts.Close()
			e := httpexpect.New(t, ts.URL)
			testCase.expectF(e.GET("/"))
		})
	}
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
			res.ServeHTTP(w, &http.Request{})
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
