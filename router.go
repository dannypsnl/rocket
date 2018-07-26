package rocket

import (
	"strings"
)

type Route struct {
	// Children route can be nil
	Children map[string]*Route
	// e.g. /user/
	// Value is `user`
	// `/` represent root route
	Value string
	// Matched means what is under the route
	// For example we can put Handler at here
	Matched *handler
}

func NewRoute(base string) *Route {
	return &Route{
		Value:    base,
		Children: make(map[string]*Route),
	}
}

func (r *Route) AddHandlerTo(route string, h *handler) {
	rs := strings.Split(route, "/")[1:]

	nexts := r.Children
	for i, rrr := range rs {
		if _, ok := nexts[rrr]; !ok {
			nexts[rrr] = NewRoute(rrr)
			if i == len(rs) - 1 {
				nexts[rrr].Matched = h
			}
		} else {
			nexts = nexts[rrr].Children
			if i == len(rs) - 1 {
				nexts[rrr].Matched = h
			}
		}
	}
}
