package rocket

import (
	"testing"

	"github.com/dannypsnl/assert"
)

func Test_route(t *testing.T) {
	assert := assert.NewTester(t)

	r := NewRoute()
	handler := &handler{route: "/world"}
	r.AddHandlerTo("/hello"+handler.route, handler)

	assert.Eq(r.Children["hello"].Children["world"].Matched.route, "/world")

	h := r.matching("/hello/world")
	assert.Eq(h.route, "/world")
}
