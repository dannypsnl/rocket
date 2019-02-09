package rocket

import (
	"fmt"
	"net/http"
)

type Route struct {
	// children route can be nil
	children map[string]*Route
	// variableRoute is prepare for route like `:name`
	variableRoute *Route
	// ownHandler means this Route has route, so not found handler would be 403(wrong method),
	// else is 404
	ownHandler bool
	// wildcardRoute is the handler of route `*path`
	wildcardRoute map[string]*handler
	// handlers stores map Method to handler of this route
	handlers map[string]*handler
	// optionsHandler stores a special handler for OPTION method handling
	optionsHandler *optionsHandler
}

func NewRoute() *Route {
	return &Route{
		children: make(map[string]*Route),
		handlers: make(map[string]*handler),
	}
}

func (r *Route) noVariableRoute() bool {
	return r.variableRoute == nil
}
func (r *Route) noWildcardRoute() bool {
	return r.wildcardRoute == nil
}

func (r *Route) addHandlerOn(baseRoute []string, h *handler) {
	// rename r to root, then following code would be more readable
	root := r
	fullRoute := append(baseRoute, h.routes...)
	currentMatchedRoute := root
	for i, r := range fullRoute {
		if isParameter(r) {
			if currentMatchedRoute.noVariableRoute() {
				currentMatchedRoute.variableRoute = NewRoute()
			}
			currentMatchedRoute = currentMatchedRoute.variableRoute
		} else if r[0] == '*' {
			h.matchedPathIndex = i - len(baseRoute)
			if currentMatchedRoute.noWildcardRoute() {
				currentMatchedRoute.wildcardRoute = make(map[string]*handler)
			}
			if _, ok := currentMatchedRoute.wildcardRoute[h.method]; !ok {
				currentMatchedRoute.addHandlerTo(currentMatchedRoute.wildcardRoute, h)
				return
			}
			panic(PanicDuplicateRoute)
		} else if _, ok := currentMatchedRoute.children[r]; !ok {
			currentMatchedRoute.children[r] = NewRoute()
			currentMatchedRoute = currentMatchedRoute.children[r]
		} else {
			currentMatchedRoute = currentMatchedRoute.children[r]
		}
	}

	if _, ok := currentMatchedRoute.handlers[h.method]; ok {
		panic(PanicDuplicateRoute)
	}
	currentMatchedRoute.addHandlerTo(currentMatchedRoute.handlers, h)
}

func (r *Route) addHandlerTo(m map[string]*handler, h *handler) {
	if r.optionsHandler == nil {
		r.optionsHandler = newOptionsHandler()
	}
	r.optionsHandler.addMethod(h.method)
	m["OPTIONS"] = r.optionsHandler.build()
	m[h.method] = h
	r.ownHandler = true
}

func (r *Route) getHandler(requestUrl []string, method string) *handler {
	if len(requestUrl) == 0 {
		if !r.ownHandler {
			return nil
		}
		if h, ok := r.handlers[method]; ok {
			return h
		}
		return newErrorHandler(http.StatusMethodNotAllowed, fmt.Sprintf(ErrorMessageForMethodNotAllowed, method))
	}

	head, restRequestUrl := requestUrl[0], requestUrl[1:]
	if router, ok := r.children[head]; ok {
		if handler := router.getHandler(restRequestUrl, method); handler != nil {
			return handler
		}
	}
	if r.variableRoute != nil {
		if handler := r.variableRoute.getHandler(restRequestUrl, method); handler != nil {
			return handler
		}
	}
	if r.wildcardRoute != nil {
		if h, hasWildcardRouteHandler := r.wildcardRoute[method]; hasWildcardRouteHandler {
			// TODO: this make handler depends on router work as its expected, should think about how to reverse their relationship
			h.addMatchedPathValueIntoContext(requestUrl...)
			return h
		}
		return newErrorHandler(http.StatusMethodNotAllowed, fmt.Sprintf(ErrorMessageForMethodNotAllowed, method))
	} else {
		return nil
	}
}

func isParameter(route string) bool {
	return route[0] == ':'
}
