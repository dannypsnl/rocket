package rocket

import (
	"testing"

	"github.com/dannypsnl/assert"
)

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
	}

	for _, testCase := range testCases {
		testMethod(t, testCase.method, testCase.handlerCreator)
	}

	t.Run("SelfDefinedMethod", func(t *testing.T) {
		route := "/"
		h := handlerByMethod(&route, func() {}, "Self")
		if h.method != "Self" {
			t.Errorf("expected: %v, actual: %v", "Self", h.method)
		}
	})
}

func testMethod(t *testing.T, method string, handlerCreator func(route string, do interface{}) *handler) {
	t.Helper()
	assert := assert.NewTester(t)
	t.Run(method, func(t *testing.T) {
		h := handlerCreator("/", func() {})
		assert.Eq(method, h.method)
	})
}
