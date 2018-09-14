package rocket

import (
	"net/http"
)

type Cookies struct {
	req *http.Request
}

func (c *Cookies) Cookie(name string) (*http.Cookie, error) {
	return c.req.Cookie(name)
}
