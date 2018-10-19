package fairing

import (
	"github.com/dannypsnl/rocket/response"
)

type ResponseHook func(*response.Response) *response.Response

type ResponseDecorator struct {
	hook ResponseHook
}

func (r *ResponseDecorator) Hook(resp *response.Response) *response.Response {
	return r.hook(resp)
}

func OnResponse(hook ResponseHook) *ResponseDecorator {
	return &ResponseDecorator{
		hook: hook,
	}
}
