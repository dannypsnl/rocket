package rocket

import (
	"testing"

	"fmt"
	"github.com/dannypsnl/assert"
)

type User struct {
	Name    string `route:"name"`
	Age     uint32 `route:"age"`
	TestInt int    `route:"test-int"`
}

func TestRoute(t *testing.T) {
	assert := assert.NewTester(t)

	route := "/world/:name/:age/:test-int"
	handler := handlerByMethod(&route, func(u *User) string {
		return fmt.Sprintf("User{name: `%s`, age: `%d`}, testInt: %d", u.Name, u.Age, u.TestInt)
	}, "GET")

	r := NewRoute()
	r.addHandlerTo("/hello"+handler.route, handler)

	t.Run("Call", func(t *testing.T) {
		actual := r.Call("/hello/world/danny/21/100")

		assert.Eq(actual, "User{name: `danny`, age: `21`}, testInt: 100")
	})
}
