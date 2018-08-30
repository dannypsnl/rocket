package rocket

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
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
	response, err := rk.handlers[r.Method].Call(r)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		response = rk.defaultHandler.Call([]reflect.Value{})[0]
	}
	switch response.(type) {
	case Html:
		w.Header().Set("Content-Type", "text/html")
	case string:
		w.Header().Set("Content-Type", "text/plain")
	}

	fmt.Fprint(w, response)
}
