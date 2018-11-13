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
