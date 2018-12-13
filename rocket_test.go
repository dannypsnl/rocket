package rocket_test

import (
	"testing"

	"github.com/dannypsnl/rocket"
)

func TestDuplicateRoutePanic(t *testing.T) {
	defer func() {
		if r := recover(); r != rocket.PanicDuplicateRoute {
			t.Error("panic message is wrong or didn't panic")
		}
	}()
	var (
		test = rocket.Get("/", func() string { return "" })
	)
	rocket.Ignite(":80888").
		Mount("/", test, test)
}
