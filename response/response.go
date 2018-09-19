package response

import (
	"net/http"
)

type Response struct {
	Headers map[string]string
	Body    interface{}
	Cookies []*http.Cookie
}

type Headers map[string]string

func New(body interface{}) *Response {
	return &Response{
		Headers: make(map[string]string),
		Body:    body,
		Cookies: make([]*http.Cookie, 0),
	}
}

func (res *Response) WithHeaders(headers Headers) *Response {
	for k, v := range headers {
		res.Headers[k] = v
	}
	return res
}

func (res *Response) WithCookies(cs ...*http.Cookie) *Response {
	for _, c := range cs {
		res.Cookies = append(res.Cookies, c)
	}
	return res
}
