package fairing

import (
	"net/http"

	"github.com/dannypsnl/rocket/response"
)

// Meta is created for fairing OnLaunch
type Meta struct {
	// Port is the serve port of rocket server
	Port string
}

// Interface specify the method that fairing could implement
type Interface interface {
	OnRequest(*http.Request) *http.Request
	OnResponse(*response.Response) *response.Response
	// OnLaunch would let you could get the metadata of rocket server
	//
	// NOTE: only work when you using `Launch()` to start server
	// won't work while you use rocket as HTTP handler
	OnLaunch(*Meta)
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
func (f Fairing) OnLaunch(*Meta) {}
