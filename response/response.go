package response

import (
	"net/http"

	"github.com/dannypsnl/rocket/cookie"
)

type Response struct {
	Headers map[string]string
	Body    interface{}
	cookies []*http.Cookie
}

type Headers map[string]string

func New(body interface{}) *Response {
	return &Response{
		Headers: make(map[string]string),
		Body:    body,
		cookies: make([]*http.Cookie, 0),
	}
}

func (res *Response) WithHeaders(headers Headers) *Response {
	for k, v := range headers {
		res.Headers[k] = v
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
