package rocket

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOnClose(t *testing.T) {
	x := 1
	onClose := func() error {
		x = 2
		return nil
	}
	err := Ignite("").OnClose(onClose).onClose()
	assert.NoError(t, err)
	assert.Equal(t, 2, x)
}
