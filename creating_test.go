package rocket

import (
	"strings"
	"testing"

	"github.com/dannypsnl/rocket/router"

	"github.com/stretchr/testify/assert"
)

type TestContext struct{}

func TestRootRouteWithUserDefinedContextWontPanic(t *testing.T) {
	if r := recover(); r != nil {
		t.Error(r)
	}
	Get("/", func(ctx *TestContext) string { return "" })
}

func TestHandlerCreatorHttpMethod(t *testing.T) {
	testCases := []struct {
		method         string
		handlerCreator func(route string, do interface{}) *handler
	}{
		{"GET", Get},
		{"POST", Post},
		{"PUT", Put},
		{"PATCH", Patch},
		{"DELETE", Delete},
		{"SelfDefined", func(route string, do interface{}) *handler {
			return handlerByMethod(&route, do, "SelfDefined")
		}},
	}

	for _, testCase := range testCases {
		testMethod(t, testCase.method, testCase.handlerCreator)
	}
}

func testMethod(t *testing.T, method string, handlerCreator func(route string, do interface{}) *handler) {
	t.Helper()
	t.Run(method, func(t *testing.T) {
		h := handlerCreator("/", func() string { return "" })
		assert.Equal(t, h.method, method)
	})
}

var (
	hello = Get("/hello/*", func() string { return "" })
)

func TestDuplicatedRoute(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Must panic when route emit duplicated!")
		}
	}()
	Ignite("").
		Mount(hello, hello)
}

func TestDuplicateRoutePanic(t *testing.T) {
	defer func() {
		if r := recover(); r != router.PanicDuplicateRoute {
			t.Error("panic message is wrong or didn't panic")
		}
	}()
	var (
		root1 = Get("/", func() string { return "" })
		root2 = Get("/", func() string { return "" })
	)
	Ignite("").
		Mount(root1, root2)
}

func TestVoidHandlingFunctionShouldBeRejected(t *testing.T) {
	defer func() {
		if r := recover(); !strings.Contains(r.(string), "handling function should be non-void function but got") {
			t.Error("panic message is wrong or didn't panic")
		}
	}()
	Get("/", func() {})
}
