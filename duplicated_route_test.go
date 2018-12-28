package rocket_test

import (
	"testing"

	"github.com/dannypsnl/rocket"
)

var hello = rocket.Get("/hello/*", func() string { return "" })

func TestDuplicatedRoute(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Must panic when route emit duplicated!")
		}
	}()
	rocket.Ignite(":8080").
		Mount("/", hello, hello)
}

func TestDuplicateRoutePanic(t *testing.T) {
	defer func() {
		if r := recover(); r != rocket.PanicDuplicateRoute {
			t.Error("panic message is wrong or didn't panic")
		}
	}()
	var (
		root1 = rocket.Get("/", func() string { return "" })
		root2 = rocket.Get("/", func() string { return "" })
	)
	rocket.Ignite(":80888").
		Mount("/", root1, root2)
}
