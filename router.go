package rocket

import (
	"fmt"
	"reflect"
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

func (r *Route) Call(url string) string {
	rs := strings.Split(url, "/")[1:]

	handler := r.matching(rs)

	contextType := handler.do.Type().In(0)
	context := reflect.New(contextType.Elem())

	hrs := strings.Split(handler.route, "/")[1:]
	for idx, route := range hrs {
		if route[0] == ':' {
			context.Elem().Field(handler.params[idx]).
				Set(reflect.ValueOf(rs[len(rs)-len(hrs)+idx]))
		}
	}

	result := handler.do.Call([]reflect.Value{context})[0]

	return fmt.Sprintf("%v", result)
}

func (r *Route) AddHandlerTo(route string, h *handler) {
	rs := strings.Split(route, "/")[1:]

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
				if route[0] == ':' {
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
