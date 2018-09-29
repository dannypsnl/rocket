package rocket

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"
)

// Rocket is our service.
type Rocket struct {
	port           string
	handlers       *Route
	defaultHandler reflect.Value
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

	handler := rk.handlers.getHandler(rs, r.Method)
	if handler != nil {
		handler.Handle(rs, w, r)
		return
	}
	// 404 Page Not Found
	w.WriteHeader(http.StatusNotFound)
	response := rk.defaultHandler.Call([]reflect.Value{})[0]
	fmt.Fprint(w, response)

}
