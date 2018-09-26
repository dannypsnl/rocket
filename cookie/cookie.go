package cookie

import (
	"net/http"
	"time"
)

type Cookie struct {
	name, value string
	path        string
	domain      string
	expires     time.Time
}

func Forget(name string) *Cookie {
	return &Cookie{
		name:    name,
		value:   "",
		path:    "/",
		expires: time.Unix(0, 0), // delete cookie is just create a cookie live time is 0
	}
}

func New(name, value string) *Cookie {
	return &Cookie{
		name:  name,
		value: value,
	}
}

func (c *Cookie) Path(path string) *Cookie {
	c.path = path
	return c
}
func (c *Cookie) Domain(domain string) *Cookie {
	c.domain = domain
	return c
}
func (c *Cookie) Expires(t time.Time) *Cookie {
	c.expires = t
	return c
}

func (c *Cookie) Generate() *http.Cookie {
	return &http.Cookie{
		Name:    c.name,
		Value:   c.value,
		Path:    c.path,
		Domain:  c.domain,
		Expires: c.expires,
	}
}
