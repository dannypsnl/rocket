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

func (r *Route) Call(req *http.Request) (interface{}, error) {
	splitRoutes := strings.Split(req.URL.Path, "/")[1:]

	rs := make([]string, 0)
	for _, rout := range splitRoutes {
		if rout != "" {
			rs = append(rs, rout)
		}
	}

	handler := r.matching(rs)
	if handler == nil {
		return nil, PageNotFound(concatString("can't found ", req.URL.Path))
	}

	return handler.do.Call(
		handler.context(rs, req),
	)[0].Interface(), nil
}

func (route *Route) addHandlerTo(routeStr string, h *handler) {
	routes := append(
		strings.Split(routeStr, "/")[1:],
		h.routes...)

	rs := make([]string, 0)
	for _, r := range routes {
		if r != "" {
			rs = append(rs, r)
		}
	}

	if len(rs) == 0 {
		route.Matched = h
		return
	}

	next := route.Children
	i := 0
	for i < len(rs) {
		r := rs[i]
		if _, ok := next[r]; !ok {
			next[r] = NewRoute()
			i++
		}
		if i != len(rs) {
			next = next[r].Children
		}
	}

	next[rs[len(rs)-1]].Matched = h
}

func (route *Route) matching(rs []string) *handler {
	if len(rs) == 0 {
		return route.Matched
	}
	useToMatch := make([]string, 0)
	next := route.Children
	i := 0
	for i < len(rs) {
		r := rs[i]
		if _, ok := next[r]; ok {
			useToMatch = append(useToMatch, r)
			i++
			if i != len(rs) {
				next = next[r].Children
			}
		} else {
			routeExist := false
			for route, _ := range next {
				if isParameter(route) {
					useToMatch = append(useToMatch, route)
					i++
					if i != len(rs) {
						next = next[route].Children
					}
					routeExist = true
					break
				} else if route[0] == '*' {
					useToMatch = append(useToMatch, route)
					next[useToMatch[len(useToMatch)-1]].Matched.addMatchedPathValueIntoContext(rs[i:]...)
					return next[useToMatch[len(useToMatch)-1]].Matched
				}
			}
			if !routeExist {
				return nil
			}
		}
	}
	return next[useToMatch[len(useToMatch)-1]].Matched
}

func isParameter(route string) bool {
	return route[0] == ':'
}
