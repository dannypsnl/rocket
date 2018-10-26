package response

import (
	"github.com/dannypsnl/assert"
	"testing"
)

func TestSetupHeaderContentType(t *testing.T) {
	testCases := []struct {
		resp        interface{}
		contentType string
	}{
		{Html(""), "text/html"},
		{Json(""), "application/json"},
		{"", "text/plain"},
	}

	for _, testCase := range testCases {
		testContentType(t, testCase.resp, testCase.contentType)
	}
}

func testContentType(t *testing.T, response interface{}, expectedContentType string) {
	t.Helper()
	assert := assert.NewTester(t)
	actualContentType := contentTypeOf(response)
	assert.Eq(actualContentType, expectedContentType)
}
