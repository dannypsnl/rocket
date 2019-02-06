package fairing

import (
	"net/http"

	"github.com/dannypsnl/rocket/response"
)

// Interface specify the method that fairing could implement
type Interface interface {
	OnRequest(*http.Request) *http.Request
	OnResponse(*response.Response) *response.Response
}

// Fairing provides default implement for your fairing by embedded it into your fairing type,
// Embedded is a good practice because fairing is designed to allowed partial implement interface
type Fairing struct{}

func (f Fairing) OnRequest(req *http.Request) *http.Request {
	return req
}
func (f Fairing) OnResponse(resp *response.Response) *response.Response {
	return resp
}
