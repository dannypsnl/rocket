package rocket

import (
	"net/http"
)

type Cookies struct {
	req *http.Request

	listOfCookie []*http.Cookie
}

func (c *Cookies) Get(name string) (*http.Cookie, error) {
	return c.req.Cookie(name)
}

func (c *Cookies) List() []*http.Cookie {
	return c.req.Cookies()
}
