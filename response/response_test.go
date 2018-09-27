package response_test

import (
	"github.com/dannypsnl/rocket/response"

	"testing"

	"net/http"

	"github.com/dannypsnl/assert"
	"github.com/dannypsnl/rocket/cookie"
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
	assert := assert.NewTester(t)

	t.Run("Cookie", func(t *testing.T) {
		res := response.New("").Cookies(
			cookie.New("c0", "v0"),
			cookie.New("c1", "v1"),
			cookie.New("c2", "v2"),
		)
		w := &fakeRespWriter{count: 0}
		res.SetCookie(w)

		assert.Eq(w.count, 3)
	})
	t.Run("Headers", func(t *testing.T) {
		res := response.New("").WithHeaders(response.Headers{
			"Content-Type": "text/plain",
		})
		w := &fakeRespWriter{count: 0}
		res.SetHeaders(w)

		assert.Eq(w.count, 1)
	})
}
