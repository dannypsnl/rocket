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
	handlers       map[string]*Route
	defaultHandler reflect.Value
}

// Mount add handler into our service.
func (rk *Rocket) Mount(route string, h *handler, hs ...*handler) *Rocket {
	verifyBase(route)

	rk.handlers[h.method].addHandlerTo(route, h)

	for _, h := range hs {
		rk.handlers[h.method].addHandlerTo(route, h)
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
	hs := make(map[string]*Route)
	// Initial internal method map
	hs["GET"] = NewRoute()
	hs["POST"] = NewRoute()
	hs["PUT"] = NewRoute()
	hs["PATCH"] = NewRoute()
	hs["DELETE"] = NewRoute()
	return &Rocket{
		port:     port,
		handlers: hs,
		defaultHandler: reflect.ValueOf(func() string {
			return "page not found"
		}),
	}
}

// ServeHTTP is prepare for http server trait, but the plan change, it need a new name.
func (rk *Rocket) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	queryIdx := strings.Index(r.URL.Path, "?")
	path := r.URL.Path
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

	handler := rk.handlers[r.Method].getHandler(rs)
	if handler != nil {
		handler.Handle(rs, w, r)
		return
	}
	// 404 Page Not Found
	w.WriteHeader(http.StatusNotFound)
	response := rk.defaultHandler.Call([]reflect.Value{})[0]
	fmt.Fprint(w, response)

}

func contentTypeOf(response interface{}) string {
	switch response.(type) {
	case Html:
		return "text/html"
	case Json:
		return "application/json"
	case string:
		return "text/plain"
	default:
		return "text/plain"
	}
}
