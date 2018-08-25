package rocket

import (
	"net/http"
	"strings"
)

type Route struct {
	// Children route can be nil
	Children map[string]*Route
	// Matched means what is under the route
	// For example we can put Handler at here
	Matched *handler
}

func NewRoute() *Route {
	return &Route{
		Children: make(map[string]*Route),
	}
}

func (r *Route) Call(req *http.Request) interface{} {
	splitRoutes := strings.Split(req.URL.Path, "/")[1:]

	rs := make([]string, 0)
	for _, rout := range splitRoutes {
		if rout != "" {
			rs = append(rs, rout)
		}
	}

	handler := r.matching(rs)

	return handler.do.Call(
		handler.context(rs, req),
	)[0].Interface()
}

func (r *Route) addHandlerTo(route string, h *handler) {
	routes := append(strings.Split(route, "/")[1:], h.routes...)

	rs := make([]string, 0)
	for _, rout := range routes {
		if rout != "" {
			rs = append(rs, rout)
		}
	}

	if len(rs) == 0 {
		r.Matched = h
		return
	}

	next := r.Children
	i := 0
	for i < len(rs) {
		rrr := rs[i]
		if _, ok := next[rrr]; !ok {
			next[rrr] = NewRoute()
			i++
		}
		if i != len(rs) {
			next = next[rrr].Children
		}
	}

	next[rs[len(rs)-1]].Matched = h
}

func (r *Route) matching(rs []string) *handler {
	if len(rs) == 0 {
		return r.Matched
	}
	useToMatch := make([]string, 0)
	next := r.Children
	i := 0
	for i < len(rs) {
		rrr := rs[i]
		if _, ok := next[rrr]; ok {
			useToMatch = append(useToMatch, rrr)
			i++
			if i != len(rs) {
				next = next[rrr].Children
			}
		} else {
			for route, _ := range next {
				if isParameter(route) {
					useToMatch = append(useToMatch, route)
					i++
					if i != len(rs) {
						next = next[route].Children
					}
					break
				}
			}
		}
	}
	return next[useToMatch[len(useToMatch)-1]].Matched
}

func isParameter(route string) bool {
	return route[0] == ':'
}
