package rocket

import (
	"fmt"
	"log"
	"net/http"
)

// Rocket is our service.
type Rocket struct {
	port     string
	handlers map[string]*Route
}

// Mount add handler into our service.
func (rk *Rocket) Mount(route string, h *handler) *Rocket {
	verifyBase(route)

	rk.handlers[h.method].addHandlerTo(route, h)

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
	hs["DELETE"] = NewRoute()
	return &Rocket{
		port:     port,
		handlers: hs,
	}
}

// ServeHTTP is prepare for http server trait, but the plan change, it need a new name.
func (rk *Rocket) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	response := rk.handlers[r.Method].Call(r.URL.Path)
	switch response.(type) {
	case Html:
		w.Header().Set("Content-Type", "text/html")
	case string:
		w.Header().Set("Content-Type", "text/plain")
	}

	fmt.Fprint(w, response)
}
