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
}

type Headers map[string]string

func New(body interface{}) *Response {
	return &Response{
		headers: make(map[string]string),
		Body:    body,
		cookies: make([]*http.Cookie, 0),
	}
}

func (res *Response) Status(code int) *Response {
	res.statusCode = code
	return res
}

func (res *Response) WithHeaders(headers Headers) *Response {
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

func (res *Response) Handle(w http.ResponseWriter) {
	w.Header().Set("Content-Type", contentTypeOf(res.Body))
	res.setHeaders(w)
	res.setCookie(w)
	res.setStatusCode(w)
	fmt.Fprint(w, res.Body)
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
