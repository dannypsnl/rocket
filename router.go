package rocket

import (
	"fmt"
	"net/http"
	"strings"
)

type Route struct {
	// Children route can be nil
	Children map[string]*Route
	// VariableRoute is prepare for route like `:name`
	VariableRoute *Route
	// OwnHandler means this Route has route, so not found handler would be 403(wrong method),
	// else is 404
	OwnHandler bool
	// PathRouteHandler is the handler of route `*path`
	PathRouteHandler map[string]*handler
	//
	Handlers map[string]*handler
	//
	optionsHandler *optionsHandler
}

func NewRoute() *Route {
	return &Route{
		Children: make(map[string]*Route),
		Handlers: make(map[string]*handler),
	}
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

	baseRouteLen := (len(rs) - len(h.routes))
	matchRoute := route
	child := route.Children
	for i, r := range rs {
		if isParameter(r) {
			if matchRoute.VariableRoute == nil {
				matchRoute.VariableRoute = NewRoute()
			}
			matchRoute = matchRoute.VariableRoute
		} else if r[0] == '*' {
			h.matchedPathIndex = i - baseRouteLen
			if matchRoute.PathRouteHandler == nil {
				matchRoute.PathRouteHandler = make(map[string]*handler)
			}
			if _, ok := matchRoute.PathRouteHandler[h.method]; !ok {
				matchRoute.addHandler(h, matchRoute.PathRouteHandler)
				return
			}
			panic(PanicDuplicateRoute)
		} else if _, ok := child[r]; !ok {
			child[r] = NewRoute()
			matchRoute = child[r]
		} else {
			matchRoute = child[r]
		}
		child = matchRoute.Children
	}

	if _, ok := matchRoute.Handlers[h.method]; ok {
		panic(PanicDuplicateRoute)
	}
	matchRoute.addHandler(h, matchRoute.Handlers)
}

func (route *Route) addHandler(h *handler, m map[string]*handler) {
	if route.optionsHandler == nil {
		route.optionsHandler = newOptionsHandler()
	}
	route.optionsHandler.addMethod(h.method)
	m["OPTIONS"] = route.optionsHandler.build()
	m[h.method] = h
	route.OwnHandler = true
}

func (route *Route) getHandler(requestUrl []string, method string) *handler {
	if len(requestUrl) == 0 {
		if !route.OwnHandler {
			return nil
		}
		if h, ok := route.Handlers[method]; ok {
			return h
		}
		return newErrorHandler(http.StatusMethodNotAllowed, fmt.Sprintf(ErrorMessageForMethodNotAllowed, method))
	}
	next := route
	for i, r := range requestUrl {
		if router, ok := next.Children[r]; ok {
			next = router
		} else if next.VariableRoute != nil {
			next = next.VariableRoute
		} else if next.PathRouteHandler != nil {
			if !route.OwnHandler {
				return nil
			}
			if h, hasPathRouteHandler := next.PathRouteHandler[method]; hasPathRouteHandler {
				// TODO: this make handler depends on router work as its expected, should think about how to reverse their relationship
				h.addMatchedPathValueIntoContext(requestUrl[i:]...)
				return h
			}
			return newErrorHandler(http.StatusMethodNotAllowed, fmt.Sprintf(ErrorMessageForMethodNotAllowed, method))
		} else {
			return nil
		}
	}
	if !next.OwnHandler {
		return nil
	}
	if h, ok := next.Handlers[method]; ok {
		return h
	}
	return newErrorHandler(http.StatusMethodNotAllowed, fmt.Sprintf(ErrorMessageForMethodNotAllowed, method))
}

func isParameter(route string) bool {
	return route[0] == ':'
}
