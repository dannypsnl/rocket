package response

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRedirect(t *testing.T) {
	res := Redirect("/")
	assert.Equal(t, "/", res.redirectPath)
}
