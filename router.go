package rocket

import (
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

func (r *Route) AddHandlerTo(route string, h *handler) {
	rs := strings.Split(route, "/")[1:]

	nexts := r.Children
	i := 0
	for i < len(rs) {
		rrr := rs[i]
		if _, ok := nexts[rrr]; !ok {
			nexts[rrr] = NewRoute()
			i++
		}
		if i != len(rs) {
			nexts = nexts[rrr].Children
		}
	}
	nexts[rs[len(rs)-1]].Matched = h
}

func (r *Route) matching(url string) *handler {
	rs := strings.Split(url, "/")[1:]

	nexts := r.Children
	i := 0
	for i < len(rs)-1 {
		rrr := rs[i]
		if _, ok := nexts[rrr]; ok {
			nexts = nexts[rrr].Children
			i++
		}
	}
	return nexts[rs[len(rs)-1]].Matched
}
