package response

import (
	"errors"
	"net/http"
	"testing"

	asserter "github.com/dannypsnl/assert"
)

func TestFileResponder(t *testing.T) {
	assert := asserter.NewTester(t)
	t.Run("Workable", func(t *testing.T) {
		f := newFileResponder("index.html")
		f.resp = New(`<h1>Title</h1>`)
		resp := f.SetContentType(ByFileNameSuffix())
		assert.Eq(resp.headers["Content-Type"], "text/html")
	})
	t.Run("Unprocessable", func(t *testing.T) {
		f := newFileResponder("")
		f.err = errors.New("")
		f.resp = New(`<h1>Title</h1>`)
		resp := f.SetContentType(ByFileNameSuffix())
		assert.Eq(resp.statusCode, http.StatusUnprocessableEntity)
	})
	t.Run("FallbackContentType", func(t *testing.T) {
		f := newFileResponder("file")
		f.resp = New(`<h1>Title</h1>`)
		r := f.SetContentType(ByFileNameSuffix())
		assert.Eq(r.headers["Content-Type"], "text/plain")
	})
}
