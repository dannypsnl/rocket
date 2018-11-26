package fairing

import (
	"net/http"
)

type RequestHook func(req *http.Request) *http.Request

type RequestDecorator struct {
	hook RequestHook
}

func (r *RequestDecorator) Invoke(req *http.Request) *http.Request {
	return r.hook(req)
}

func OnRequest(hook RequestHook) *RequestDecorator {
	return &RequestDecorator{
		hook: hook,
	}
}
