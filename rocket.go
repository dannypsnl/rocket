package rocket

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Rocket struct {
	port     string
	matchs   []string
	handlers map[string]Handler
}

func (r *Rocket) Mount(route string, h Handler) *Rocket {
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
		if r == ':' {
			if firstTime {
				match = route[:i-1]
				firstTime = false
			}
			start = i + 1
			open = true
		}
		if r == '*' {
			match = route[:i-1]
			params = append(params, route[i+1:])
			break
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
	r.matchs = append(r.matchs, match)
	r.handlers[match] = h
	return r
}

func (rk *Rocket) Launch() {
	server := &http.Server{
		Addr:         rk.port,
		Handler:      rk,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	server.ListenAndServe()
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
	for _, m := range rk.matchs {
		if strings.HasPrefix(r.URL.Path, m) {
			match = m
			paramsPart = r.URL.Path[len(m):]
			break
		}
	}
	var params []string
	for _, param := range strings.Split(paramsPart, "/") {
		if strings.Compare(param, "") != 0 {
			params = append(params, param)
		}
	}
	h := rk.handlers[match]

	Context := make(map[string]string)
	for i, param := range h.Params {
		Context[param] = params[i]
	}
	fmt.Fprintf(w, h.Do(Context))
}
