package response

import (
	"net/http"
	"testing"

	asserter "github.com/dannypsnl/assert"
)

func TestFile(t *testing.T) {
	assert := asserter.NewTester(t)
	t.Run("ExistFile", func(t *testing.T) {
		response := File("../test_data/test.html")
		assert.Eq(response.headers["Content-Type"], "text/html")
	})
	t.Run("NoneExistFile", func(t *testing.T) {
		response := File("test_data/test.html")
		assert.Eq(response.statusCode, http.StatusNotFound)
	})
}
