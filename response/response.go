package response

import (
	"fmt"
	"net/http"

	"github.com/dannypsnl/rocket/cookie"
)

// Response provide an abstraction for detailed HTTP response
type Response struct {
	headers map[string]string
	Body    interface{}

	cookies    []*http.Cookie
	statusCode int

	keepFunc func(w http.ResponseWriter)
}

// Headers helps code be more readable
type Headers map[string]string

// New would create a new response by provided body
func New(body interface{}) *Response {
	return (&Response{
		headers: make(map[string]string),
		Body:    body,
		cookies: make([]*http.Cookie, 0),
	}).ContentType(contentTypeOf(body))
}

// isValidStatusCode would return a status code is in valid range or not
//
// According to https://tools.ietf.org/html/rfc2616#section-6.1.1
// > The Status-Code element is a 3-digit integer
//
// implements this check
func isValidStatusCode(code int) bool {
	return code > 99 && code < 1000
}

// Status would change the status code of response by provided code,
// it would panic if you call it twice on the same response or you provide a invalid status code.
//
// NOTE:
// According to https://tools.ietf.org/html/rfc2616#section-6.1.1
// > The Status-Code element is a 3-digit integer
func (res *Response) Status(code int) *Response {
	if !isValidStatusCode(code) {
		panic(fmt.Errorf("reject invalid status code: %d", code))
	}
	if res.statusCode != 0 {
		panic("status code already been set")
	}
	res.statusCode = code
	return res
}

// ContentType would change content-type of response by provided value
func (res *Response) ContentType(value string) *Response {
	return res.Headers(Headers{
		"Content-Type": value,
	})
}

// Headers would update headers of response by provided headers
func (res *Response) Headers(headers Headers) *Response {
	for k, v := range headers {
		res.headers[k] = v
	}
	return res
}

// Cookies would update cookies of response by provided cookies
func (res *Response) Cookies(cookies ...*cookie.Cookie) *Response {
	for _, c := range cookies {
		res.cookies = append(res.cookies, c.Generate())
	}
	return res
}

func (res *Response) keep(keepFunc func(w http.ResponseWriter)) *Response {
	res.keepFunc = keepFunc
	return res
}

// WriteTo is not exported for you, so do not call it directly
func (res *Response) WriteTo(w http.ResponseWriter) {
	res.setHeaders(w)
	res.setCookie(w)
	res.setStatusCode(w)
	fmt.Fprint(w, res.Body)
	if res.keepFunc != nil {
		res.keepFunc(w)
	}
}

func (res *Response) setHeaders(w http.ResponseWriter) {
	for k, v := range res.headers {
		w.Header().Set(k, v)
	}
}

func (res *Response) setCookie(w http.ResponseWriter) {
	for _, c := range res.cookies {
		http.SetCookie(w, c)
	}
}

func (res *Response) setStatusCode(w http.ResponseWriter) {
	if res.statusCode != 0 {
		w.WriteHeader(res.statusCode)
	}
}

func Stream(f func(http.ResponseWriter)) *Response {
	return New("").
		keep(f)
}
