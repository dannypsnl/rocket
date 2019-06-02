package cookie

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCookie(t *testing.T) {
	t.Run("Domain", func(t *testing.T) {
		c := New("test", "value").
			Domain("example.com")
		assert.Equal(t, "example.com", c.domain)
	})
	t.Run("MaxAge", func(t *testing.T) {
		c := New("test", "value").MaxAge(10)
		assert.Equal(t, 10, c.maxAge)
	})
	t.Run("Forget", func(t *testing.T) {
		c := Forget("test")
		assert.Equal(t, "test", c.name)
		assert.Equal(t, "/", c.path)
		assert.Equal(t, time.Unix(0, 0), c.expires)
	})
	t.Run("Generate", func(t *testing.T) {
		realCookie := New("test", "value").
			Generate()
		assert.Equal(t, "test", realCookie.Name)
	})
}
