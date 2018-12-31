package fairing

import (
	"net/http"

	"github.com/dannypsnl/rocket/response"
)

type FairingInterface interface {
	OnRequest(*http.Request) *http.Request
	OnResponse(*response.Response) *response.Response
}

type Fairing struct{}

func (f Fairing) OnRequest(req *http.Request) *http.Request {
	return req
}
func (f Fairing) OnResponse(resp *response.Response) *response.Response {
	return resp
}
