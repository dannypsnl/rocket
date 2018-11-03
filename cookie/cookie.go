package cookie

import (
	"net/http"
	"time"
)

type Cookie struct {
	name, value string
	path        string
	domain      string
	// maxAge=0 means no 'Max-Age' attribute specified
	// maxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
	// maxAge>0 means Max-Age attribute present and given in seconds
	maxAge  int
	expires time.Time
}

// Forget would delete the cookie that name is provided name
func Forget(name string) *Cookie {
	return &Cookie{
		name:  name,
		value: "",
		path:  "/",
		// although maxAge: -1 can delete cookie too, but not all platform can recognize it, so use expires at here
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

// MaxAge is set function for maxAge of Cookie
//
// * maxAge=0 means no 'Max-Age' attribute specified
//
// * maxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
//
// * maxAge>0 means Max-Age attribute present and given in seconds
func (c *Cookie) MaxAge(maxAge int) *Cookie {
	c.maxAge = maxAge
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
