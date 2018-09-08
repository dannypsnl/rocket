package rocket

import (
	"net/http"
	"strings"
)

type Route struct {
	// Children route can be nil
	Children map[string]*Route
	// VariableRoute is prepare for route like `:name`
	VariableRoute *Route
	// PathRouteHandler is the handler of route `*path`
	PathRouteHandler *handler
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
		if isParameter(r) {
			if matchRoute.VariableRoute == nil {
				matchRoute.VariableRoute = NewRoute()
			}
			matchRoute = matchRoute.VariableRoute
		} else if r[0] == '*' {
			if matchRoute.PathRouteHandler == nil {
				matchRoute.PathRouteHandler = h
				return
			}
			panic("Duplicated route")
		} else if _, ok := child[r]; !ok {
			child[r] = NewRoute()
			matchRoute = child[r]
		} else {
			matchRoute = child[r]
		}
		child = matchRoute.Children
	}

	matchRoute.Matched = h
}

func (route *Route) matching(requestUrl []string) *handler {
	if len(requestUrl) == 0 {
		return route.Matched
	}
	next := route
	for i, r := range requestUrl {
		if router, ok := next.Children[r]; ok {
			next = router
		} else if next.VariableRoute != nil {
			next = next.VariableRoute
		} else if next.PathRouteHandler != nil {
			next.PathRouteHandler.addMatchedPathValueIntoContext(requestUrl[i:]...)
			return next.PathRouteHandler
		} else {
			return nil
		}
	}
	return next.Matched
}

func isParameter(route string) bool {
	return route[0] == ':'
}
