package rocket_test

import (
	"net/http/httptest"
	"testing"

	"github.com/dannypsnl/rocket"

	"github.com/gavv/httpexpect"
)

var (
	forTestHandler = rocket.Get("/", func() string { return "" })
)

func TestOptionsMethod(t *testing.T) {
	rk := rocket.Ignite(":8081").
		Mount("/", forTestHandler)
	ts := httptest.NewServer(rk)
	defer ts.Close()
	e := httpexpect.New(t, ts.URL)

	e.OPTIONS("/").
		Expect().
		Header("Allow").
		Equal("OPTIONS, GET")
}
