package cookie

import (
	"testing"

	asserter "github.com/dannypsnl/assert"
)

func TestCookie(t *testing.T) {
	assert := asserter.NewTester(t)
	t.Run("MaxAge", func(t *testing.T) {
		c := New("test", "value").MaxAge(10)
		assert.Eq(10, c.maxAge)
		c = New("test", "value")
		assert.Eq(0, c.maxAge)
	})
}
