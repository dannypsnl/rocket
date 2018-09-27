package rocket_test

import (
	"github.com/dannypsnl/rocket"

	"testing"
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
