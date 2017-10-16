package rocket

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Rocket struct {
	port     string
	matchs   []string
	handlers map[string]Handler
}

func (r *Rocket) Mount(route string, h Handler) *Rocket {
	route += h.Route
	// '/:id' is params in url.
	// '/*filepath' is params about filepath.
	// '/home, data' is params from post method.
	match, params := split(route)
	h.params = params
	r.matchs = append(r.matchs, match)
	r.handlers[match] = h
	return r
}

func (rk *Rocket) Launch() {
	http.HandleFunc("/", rk.ServeHTTP)
	log.Fatal(http.ListenAndServe(rk.port, nil))
}

func Ignite(port string) *Rocket {
	return &Rocket{
		port:     port,
		handlers: make(map[string]Handler),
	}
}

func (rk *Rocket) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var match string
	var paramsPart string
	for _, m := range rk.matchs { // rk.matchs are those static routes
		if strings.HasPrefix(r.URL.Path, m) {
			match = m
			paramsPart = r.URL.Path[len(m):]
			break
		}
	}
	h := rk.handlers[match]
	Context := make(map[string]string)
	if strings.HasPrefix(h.params[0], "*") {
		Context[h.params[0][1:]] = paramsPart
	} else {
		var params []string
		for _, param := range strings.Split(paramsPart, "/") {
			if strings.Compare(param, "") != 0 {
				params = append(params, param)
			}
		}
		for i, param := range h.params {
			Context[param] = params[i]
		}
	}

	fmt.Fprintf(w, h.Do(Context))
}
