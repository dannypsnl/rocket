package response

import (
	asserter "github.com/dannypsnl/assert"
	"testing"
)

func TestFileResponser(t *testing.T) {
	assert := asserter.NewTester(t)
	f := newFileResponser("index.html")
	f.resp = New(`<h1>Title</h1>`)
	resp := f.ByFileSuffix(DefaultContentTypes)
	assert.Eq(resp.headers["Content-Type"], "text/html")
}
