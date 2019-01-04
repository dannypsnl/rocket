package rocket

import (
	"testing"

	"github.com/dannypsnl/assert"
)

type (
	TestContext struct {
	}
)

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
	assert := assert.NewTester(t)
	t.Run(method, func(t *testing.T) {
		h := handlerCreator("/", func() {})
		assert.Eq(method, h.method)
	})
}
