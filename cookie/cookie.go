package cookie

import (
	"net/http"
	"time"
)

// Cookie provides an abstraction of cookie in rocket
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
	return New(name, "").Path("/").
		// although maxAge: -1 can delete cookie too, but not all platform can recognize it, so use expires at here
		Expires(time.Unix(0, 0))
}

// New create a new cookie
func New(name, value string) *Cookie {
	return &Cookie{
		name:  name,
		value: value,
	}
}

// Path set path of cookie
func (c *Cookie) Path(path string) *Cookie {
	c.path = path
	return c
}

// Domain set domain of cookie
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

// Expires set expires of cookie
func (c *Cookie) Expires(time time.Time) *Cookie {
	c.expires = time
	return c
}

// Generate is not prepare for you
func (c *Cookie) Generate() *http.Cookie {
	return &http.Cookie{
		Name:    c.name,
		Value:   c.value,
		Path:    c.path,
		Domain:  c.domain,
		MaxAge:  c.maxAge,
		Expires: c.expires,
	}
}
