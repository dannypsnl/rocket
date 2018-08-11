package rocket

import (
	"reflect"
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
	r.AddHandlerTo("/hello"+handler.route, handler)

	h := r.matching("/hello/world/0")

	u := &User{Id: "0"}
	result := h.do.Call([]reflect.Value{
		reflect.ValueOf(u),
	})[0]

	assert.Eq(result.Interface(), "0")
}
