package rocket

import (
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

	rk.handlers[h.method].AddHandlerTo(route+h.route, h)

	return rk
}

// Launch shoot our service.(start server)
func (rk *Rocket) Launch() {
	http.HandleFunc("/", rk.serveLoop)
	log.Fatal(http.ListenAndServe(rk.port, nil))
}

// Ignite initial service by port.
func Ignite(port string) *Rocket {
	hs := make(map[string]*Route)
	// Initial internal method map
	hs["GET"] = &Route{}
	hs["POST"] = &Route{}
	hs["PUT"] = &Route{}
	hs["DELETE"] = &Route{}
	return &Rocket{
		port:     port,
		handlers: hs,
	}
}

// serveLoop is prepare for http server trait, but the plan change, it need a new name.
func (rk *Rocket) serveLoop(w http.ResponseWriter, r *http.Request) {
	//h := rk.handlers[r.Method].matching(r.URL)
}
