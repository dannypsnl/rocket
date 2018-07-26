package rocket

import (
	"log"
	"net/http"
	"strings"
)

// Rocket is our service.
type Rocket struct {
	port     string
	handlers map[string]*Route
}

// Mount add handler into our service.
func (rk *Rocket) Mount(route string, h *handler) *Rocket {
	verifyBase(route)

	root := rk.handlers[h.method]
	nexts := root.Children
	for _, r := range strings.Split(route+h.route, "/") {
		if nexts == nil {
			r := &Route{Value: r}
			nexts = []*Route{r}
		} else {
			for _, next := range nexts {
				if next.Value == r {
					nexts = next.Children
				}
			}
		}
	}
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
	hs["GET"] = &Route{Value: "/"}
	hs["POST"] = &Route{Value: "/"}
	hs["PUT"] = &Route{Value: "/"}
	hs["DELETE"] = &Route{Value: "/"}
	return &Rocket{
		port:     port,
		handlers: hs,
	}
}

// serveLoop is prepare for http server trait, but the plan change, it need a new name.
func (rk *Rocket) serveLoop(w http.ResponseWriter, r *http.Request) {
}
