package rocket

import (
	"github.com/dannypsnl/assert"
	"testing"

	"fmt"
)

type User struct {
	Name string `route:"name"`
	Age  uint32 `route:"age"`
}

func TestRoute(t *testing.T) {
	assert := assert.NewTester(t)

	route := "/world/:name/:age"
	handler := handlerByMethod(&route, func(u *User) string {
		return fmt.Sprintf("User{name: `%s`, age: `%d`}", u.Name, u.Age)
	}, "GET")

	r := NewRoute()
	r.addHandlerTo("/hello", handler)

	t.Run("Call", func(t *testing.T) {
		actual := r.Call("/hello/world/danny/21")

		assert.Eq(actual, "User{name: `danny`, age: `21`}")
	})
}
