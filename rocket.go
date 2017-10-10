package rocket

import (
	"fmt"
	"log"
	"net/http"

	"rocket/routes"
)

type Rocket struct {
	port     string
	handlers map[string]routes.Handler
}

func (r *Rocket) Mount(route string, h routes.Handler) *Rocket {
	route += h.Route
	match := ""

	firstTime := true
	open := false
	start := 0
	var params []string
	// '/:id' is params in url.
	// '/*filepath' is params about filepath.
	// '/home, data' is params from post method.
	for i, r := range route {
		if r == ':' || r == '*' {
			if firstTime {
				match = route[:i-1]
				firstTime = false
			}
			start = i + 1
			open = true
		}
		if i == len(route)-1 {
			params = append(params, route[start:i+1])
		}
		if r == '/' && open {
			// Get param setting string.
			params = append(params, route[start:i])
			open = false
		}
	}
	h.Params = params
	r.handlers[match] = h
	return r
}

func (rk *Rocket) Launch() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		h := rk.handlers[r.URL.Path]
		fmt.Fprintf(w, "%v", h.Params)
		fmt.Fprintf(w, h.Do())
	})
	log.Fatal(http.ListenAndServe(rk.port, nil))
}

func Ignite(port string) *Rocket {
	return &Rocket{
		port:     port,
		handlers: make(map[string]routes.Handler),
	}
}
