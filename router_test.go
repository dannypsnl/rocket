package rocket

import (
	"testing"

	"github.com/dannypsnl/assert"
)

type User struct {
	Id string `route:"id"`
}

func TestRoute(t *testing.T) {
	assert := assert.NewTester(t)

	route := "/world/:id"
	handler := handlerByMethod(&route, func(u *User) string {
		return u.Id
	}, "GET")

	r := NewRoute()
	r.addHandlerTo("/hello"+handler.route, handler)

	t.Run("Call", func(t *testing.T) {
		actual := r.Call("/hello/world/0")

		assert.Eq(actual, "0")
	})
}
