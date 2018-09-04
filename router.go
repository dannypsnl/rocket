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
	queryIdx := strings.Index(req.URL.Path, "?")
	path := req.URL.Path
	if queryIdx > -1 {
		path = path[:queryIdx]
	}
	splitRoutes := strings.Split(path, "/")[1:]

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
	} else if len(rs) == 1 {
		route.Children[rs[0]].Matched = h
		return
	}

	next := route.Children
	i := 0
	for i < len(rs) {
		r := rs[i]
		if _, ok := next[r]; !ok {
			next[r] = NewRoute()
		}
		// increase i whether create new route or not
		i++
		if i != len(rs) {
			next = next[r].Children
		}
	}

	next[rs[len(rs)-1]].Matched = h
}

func (route *Route) matching(requestUrl []string) *handler {
	if len(requestUrl) == 0 {
		return route.Matched
	}
	next := route.Children
	i := 0
	for i < len(requestUrl) {
		r := requestUrl[i]
		if router, ok := next[r]; ok {
			i++
			if i != len(requestUrl) {
				next = next[r].Children
			} else {
				return router.Matched
			}
		} else {
			found := false
			for route, _ := range next {
				if isParameter(route) {
					i++
					if i != len(requestUrl) {
						found = true
						next = next[route].Children
					} else {
						return next[route].Matched
					}
				} else if route[0] == '*' {
					next[route].Matched.addMatchedPathValueIntoContext(requestUrl[i:]...)
					return next[route].Matched
				}
			}
			if !found {
				return nil
			}
		}
	}
	return nil
}

func isParameter(route string) bool {
	return route[0] == ':'
}
