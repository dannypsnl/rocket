package response

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFile(t *testing.T) {
	t.Run("ExistFile", func(t *testing.T) {
		response := File("../test_data/test.html")
		assert.Equal(t, "text/html", response.headers["Content-Type"])
	})
	t.Run("NoneExistFile", func(t *testing.T) {
		response := File("test_data/test.html")
		assert.Equal(t, http.StatusNotFound, response.statusCode)
	})
}
