package fairing

import (
	"github.com/dannypsnl/rocket/response"
)

type ResponseHook func(*response.Response) *response.Response

type Response struct {
	hook ResponseHook
}

func (r *Response) Hook(resp *response.Response) *response.Response {
	return r.hook(resp)
}

func OnResponse(hook ResponseHook) *Response {
	return &Response{
		hook: hook,
	}
}
