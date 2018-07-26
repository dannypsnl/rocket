package rocket

import (
	"testing"
	"github.com/dannypsnl/assert"
)

func Test_route(t *testing.T) {
	assert := assert.NewTester(t)

	r := NewRoute("/")
	handler := &handler{route: "/world"}
	r.AddHandlerTo("/hello", handler)

	assert.Eq(len(r.Children), 1)
}
