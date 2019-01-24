package response

import (
	"fmt"
	"net/http"

	"github.com/dannypsnl/rocket/cookie"
)

type Response struct {
	headers map[string]string
	Body    interface{}

	cookies    []*http.Cookie
	statusCode int

	keepFunc func(w http.ResponseWriter)
}

type Headers map[string]string

func New(body interface{}) *Response {
	return &Response{
		headers: map[string]string{
			"Content-Type": contentTypeOf(body),
		},
		Body:    body,
		cookies: make([]*http.Cookie, 0),
	}
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

func (res *Response) Headers(headers Headers) *Response {
	for k, v := range headers {
		res.headers[k] = v
	}
	return res
}

func (res *Response) Cookies(cs ...*cookie.Cookie) *Response {
	for _, c := range cs {
		res.cookies = append(res.cookies, c.Generate())
	}
	return res
}

func (res *Response) keep(keepFunc func(w http.ResponseWriter)) *Response {
	res.keepFunc = keepFunc
	return res
}

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

func Stream(f func(chan<- []byte, func())) *Response {
	return New("").
		keep(func(w http.ResponseWriter) {
			ch := make(chan []byte)
			stopCh := make(chan struct{})
			done := func() {
				stopCh <- struct{}{}
			}
			if f == nil {
				return
			}
			go f(ch, done)
			for {
				select {
				case data := <-ch:
					fmt.Fprint(w, data)
				case <-stopCh:
					return
				}
			}
		})
}
