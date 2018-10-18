package response

import (
	"errors"
	asserter "github.com/dannypsnl/assert"
	"net/http"
	"testing"
)

func TestFileResponser(t *testing.T) {
	assert := asserter.NewTester(t)
	t.Run("Workable", func(t *testing.T) {
		f := newFileResponser("index.html")
		f.resp = New(`<h1>Title</h1>`)
		resp := f.ByFileSuffix(DefaultContentTypes)
		assert.Eq(resp.headers["Content-Type"], "text/html")
	})
	t.Run("Unprocessable", func(t *testing.T) {
		f := newFileResponser("")
		f.err = errors.New("")
		f.resp = New(`<h1>Title</h1>`)
		resp := f.ByFileSuffix(DefaultContentTypes)
		assert.Eq(resp.statusCode, http.StatusUnprocessableEntity)
	})
	t.Run("FallbackContentType", func(t *testing.T) {
		f := newFileResponser("file")
		f.resp = New(`<h1>Title</h1>`)
		r := f.ByFileSuffix(DefaultContentTypes)
		assert.Eq(r.headers["Content-Type"], "text/plain")
	})
}
