package cookie

import (
	"testing"
	"time"

	asserter "github.com/dannypsnl/assert"
)

func TestCookie(t *testing.T) {
	assert := asserter.NewTester(t)
	t.Run("Domain", func(t *testing.T) {
		c := New("test", "value").
			Domain("example.com")
		assert.Eq(c.domain, "example.com")
	})
	t.Run("Generate", func(t *testing.T) {
		realCookie := New("test", "value").
			Generate()
		assert.Eq(realCookie.Name, "test")
	})
	t.Run("MaxAge", func(t *testing.T) {
		c := New("test", "value").MaxAge(10)
		assert.Eq(c.maxAge, 10)
		c = New("test", "value")
		assert.Eq(c.maxAge, 0)
	})
	t.Run("Forget", func(t *testing.T) {
		c := Forget("test")
		assert.Eq(c.name, "test")
		assert.Eq(c.path, "/")
		assert.Eq(c.expires, time.Unix(0, 0))
	})
}
