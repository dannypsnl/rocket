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
	port           string
	handlers       *Route
	defaultHandler reflect.Value
	responseHook   *fairing.ResponseDecorator
}

// Mount add handler into our service.
func (rk *Rocket) Mount(route string, h *handler, hs ...*handler) *Rocket {
	verifyBase(route)

	rk.handlers.addHandlerTo(route, h)

	for _, h := range hs {
		rk.handlers.addHandlerTo(route, h)
	}

	return rk
}

// Attach add fairing to lifecycle of each request to response
func (rk *Rocket) Attach(f interface{}) *Rocket {
	switch v := f.(type) {
	case *fairing.ResponseDecorator:
		rk.responseHook = v
	default:
		panic("not support fairing")
	}
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
		port:     port,
		handlers: NewRoute(),
		defaultHandler: reflect.ValueOf(func() string {
			return "page not found"
		}),
	}
}

// ServeHTTP is prepare for http server trait, but the plan change, it need a new name.
func (rk *Rocket) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	queryIdx := strings.Index(r.URL.Path, "?")
	path := r.URL.Path
	// has query string
	if queryIdx > -1 {
		path = path[:queryIdx]
	}

	splitRoutes := strings.Split(path, "/")[1:]

	rs := make([]string, 0)
	for _, rout := range splitRoutes {
		if rout != "" {
			rs = append(rs, rout)
		}
	}

	// get response
	handler := rk.handlers.getHandler(rs, r.Method)
	var resp *response.Response
	if handler != nil {
		resp = handler.Handle(rs, r)
	} else {
		body := rk.defaultHandler.Call([]reflect.Value{})[0]
		resp = response.New(body).Status(http.StatusNotFound)
	}

	if rk.responseHook != nil {
		resp = rk.responseHook.Hook(resp)
	}
	resp.Handle(w)
}
