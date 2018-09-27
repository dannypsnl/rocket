package response

import (
	"net/http"

	"github.com/dannypsnl/rocket/cookie"
)

type Response struct {
	headers map[string]string
	Body    interface{}
	cookies []*http.Cookie
}

type Headers map[string]string

func New(body interface{}) *Response {
	return &Response{
		headers: make(map[string]string),
		Body:    body,
		cookies: make([]*http.Cookie, 0),
	}
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

func (res *Response) SetCookie(w http.ResponseWriter) {
	for _, c := range res.cookies {
		http.SetCookie(w, c)
	}
}

func (res *Response) SetHeaders(w http.ResponseWriter) {
	for k, v := range res.headers {
		w.Header().Set(k, v)
	}
}
