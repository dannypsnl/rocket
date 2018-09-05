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

	matchRoute := route
	child := route.Children
	for _, r := range rs {
		// create route if child[r] is nil
		if _, ok := child[r]; !ok {
			child[r] = NewRoute()
		}
		matchRoute = child[r]
		child = matchRoute.Children
	}

	matchRoute.Matched = h
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
