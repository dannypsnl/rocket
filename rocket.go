package rocket

import (
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/dannypsnl/rocket/fairing"
	"github.com/dannypsnl/rocket/response"
)

// Rocket is our service.
type Rocket struct {
	port          string
	handlers      *Route
	listOfFairing []fairing.FairingInterface

	defaultHandler reflect.Value
	defaultResp    *response.Response
}

// Mount add handler into our service.
func (rk *Rocket) Mount(routeStr string, h *handler, hs ...*handler) *Rocket {
	verifyBase(routeStr)

	route := make([]string, 0)
	for _, r := range strings.Split(routeStr, "/")[1:] {
		if r != "" {
			route = append(route, r)
		}
	}
	rk.handlers.addHandlerOn(route, h)

	for _, h := range hs {
		rk.handlers.addHandlerOn(route, h)
	}

	return rk
}

// Attach add fairing to lifecycle of each request to response
func (rk *Rocket) Attach(f fairing.FairingInterface) *Rocket {
	rk.listOfFairing = append(rk.listOfFairing, f)
	return rk
}

func (rk *Rocket) Default(do interface{}) *Rocket {
	rk.defaultHandler = reflect.ValueOf(do)
	return rk
}

// Launch shoot our service.(start server)
func (rk *Rocket) Launch() {
	http.HandleFunc("/", rk.ServeHTTP)
	log.Fatal(http.ListenAndServe(rk.port, nil))
}

// Ignite initial service by port.
func Ignite(port string) *Rocket {
	return &Rocket{
		port:          port,
		handlers:      NewRoute(),
		listOfFairing: make([]fairing.FairingInterface, 0),
		defaultHandler: reflect.ValueOf(func() string {
			return "page not found"
		}),
	}
}

// ServeHTTP is prepare for http server trait, but the plan change, it need a new name.
func (rk *Rocket) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	splitRoutes := strings.Split(r.URL.Path, "/")[1:]

	reqURL := make([]string, 0)
	for _, rout := range splitRoutes {
		if rout != "" {
			reqURL = append(reqURL, rout)
		}
	}

	for _, f := range rk.listOfFairing {
		f.OnRequest(r)
	}

	// get response
	handler := rk.handlers.getHandler(reqURL, r.Method)
	var resp *response.Response
	if handler != nil {
		resp = handler.Handle(reqURL, r)
	} else {
		resp = rk.defaultResponse()
	}

	for _, f := range rk.listOfFairing {
		resp = f.OnResponse(resp)
	}
	resp.WriteTo(w)
}

func (rk *Rocket) defaultResponse() *response.Response {
	if rk.defaultResp != nil {
		return rk.defaultResp
	}
	rk.defaultResp = response.New(
		rk.defaultHandler.Call([]reflect.Value{})[0],
	).Status(http.StatusNotFound)
	return rk.defaultResp
}
