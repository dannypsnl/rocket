package response

import (
	"testing"

	"github.com/dannypsnl/rocket/response/content_type"

	"github.com/stretchr/testify/assert"
)

func TestSetupHeaderContentType(t *testing.T) {
	testCases := []struct {
		resp        interface{}
		contentType string
	}{
		{Html(""), content_type.HTML},
		{Json(""), content_type.JSON},
		{"", content_type.TextPlain},
	}

	for _, testCase := range testCases {
		assertContentType(t, testCase.resp, testCase.contentType)
	}
}

func assertContentType(t *testing.T, response interface{}, expectedContentType string) {
	t.Helper()
	actualContentType := contentTypeOf(response)
	assert.Equal(t, expectedContentType, actualContentType)
}
