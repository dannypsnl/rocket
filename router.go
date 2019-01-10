package rocket

import (
	"fmt"
	"net/http"
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

func (route *Route) noVariableRoute() bool {
	return route.VariableRoute == nil
}
func (route *Route) noWildcardRoute() bool {
	return route.PathRouteHandler == nil
}

func (route *Route) addHandlerOn(baseRoute []string, h *handler) {
	// rename route to root, then following code would be more readable
	root := route
	fullRoute := append(baseRoute, h.routes...)
	currentMatchedRoute := root
	for i, r := range fullRoute {
		if isParameter(r) {
			if currentMatchedRoute.noVariableRoute() {
				currentMatchedRoute.VariableRoute = NewRoute()
			}
			currentMatchedRoute = currentMatchedRoute.VariableRoute
		} else if r[0] == '*' {
			h.matchedPathIndex = i - len(baseRoute)
			if currentMatchedRoute.noWildcardRoute() {
				currentMatchedRoute.PathRouteHandler = make(map[string]*handler)
			}
			if _, ok := currentMatchedRoute.PathRouteHandler[h.method]; !ok {
				currentMatchedRoute.addHandlerTo(currentMatchedRoute.PathRouteHandler, h)
				return
			}
			panic(PanicDuplicateRoute)
		} else if _, ok := currentMatchedRoute.Children[r]; !ok {
			currentMatchedRoute.Children[r] = NewRoute()
			currentMatchedRoute = currentMatchedRoute.Children[r]
		} else {
			currentMatchedRoute = currentMatchedRoute.Children[r]
		}
	}

	if _, ok := currentMatchedRoute.Handlers[h.method]; ok {
		panic(PanicDuplicateRoute)
	}
	currentMatchedRoute.addHandlerTo(currentMatchedRoute.Handlers, h)
}

func (route *Route) addHandlerTo(m map[string]*handler, h *handler) {
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
