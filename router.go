package rocket

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/dannypsnl/rocket/response"
)

type Route struct {
	// children route can be nil
	children map[string]*Route
	// variableRoute is prepare for route like `:name`
	variableRoute *Route
	// ownHandler means this Route has route, so not found handler would be 403(wrong method),
	// else is 404
	ownHandler bool
	// wildcardMethodHandlers is the handler of route `*path`
	wildcardMethodHandlers map[string]*handler
	// methodHandlers stores map Method to handler of this route
	methodHandlers map[string]*handler
	// optionsHandler stores a special handler for OPTION method handling
	optionsHandler *optionsHandler
}

func NewRoute() *Route {
	return &Route{
		children:       make(map[string]*Route),
		methodHandlers: make(map[string]*handler),
	}
}

func (r *Route) mustGetVariableRoute() *Route {
	if r.variableRoute == nil {
		r.variableRoute = NewRoute()
	}
	return r.variableRoute
}
func (r *Route) prepareWildcardRoute() {
	if r.wildcardMethodHandlers == nil {
		r.wildcardMethodHandlers = make(map[string]*handler)
	}
}

const PanicDuplicateRoute = "Duplicate Route"

func (root *Route) addHandler(h *handler) {
	fullRoute := h.routes
	curRoute := root
	for i, r := range fullRoute {
		if r[0] == ':' {
			curRoute = curRoute.mustGetVariableRoute()
			continue
		} else if r[0] == '*' {
			h.matchedPathIndex = i
			curRoute.prepareWildcardRoute()
			if _, ok := curRoute.wildcardMethodHandlers[h.method]; ok {
				panic(PanicDuplicateRoute)
			}
			curRoute.addHandlerOn(curRoute.wildcardMethodHandlers, h)
			return
		} else if _, ok := curRoute.children[r]; !ok {
			curRoute.children[r] = NewRoute()
		}
		curRoute = curRoute.children[r]
	}

	if _, ok := curRoute.methodHandlers[h.method]; ok {
		panic(PanicDuplicateRoute)
	}
	curRoute.addHandlerOn(curRoute.methodHandlers, h)
}

func (r *Route) addHandlerOn(m map[string]*handler, h *handler) {
	if r.optionsHandler == nil {
		r.optionsHandler = newOptionsHandler()
	}
	r.optionsHandler.addMethod(h.method)
	m["OPTIONS"] = r.optionsHandler.build()
	m[h.method] = h
	r.ownHandler = true
}

const ErrorMessageForMethodNotAllowed = "request resource does not support http method '%s'"

func (r *Route) getHandler(requestUrl []string, method string) *handler {
	if len(requestUrl) == 0 {
		if !r.ownHandler {
			return nil
		}
		if h, ok := r.methodHandlers[method]; ok {
			return h
		}
		return responseNotAllowed(method)
	}

	headOfUrl, onRestUrl := requestUrl[0], requestUrl[1:]
	if router, ok := r.children[headOfUrl]; ok {
		if h := router.getHandler(onRestUrl, method); h != nil {
			return h
		}
	}
	if r.variableRoute != nil {
		if h := r.variableRoute.getHandler(onRestUrl, method); h != nil {
			return h
		}
	}
	if r.wildcardMethodHandlers != nil {
		if h, hasWildcardRouteHandler := r.wildcardMethodHandlers[method]; hasWildcardRouteHandler {
			return h
		}
		return responseNotAllowed(method)
	} else {
		return nil
	}
}

func responseNotAllowed(method string) *handler {
	return newHandler(reflect.ValueOf(func() *response.Response {
		return response.New(fmt.Sprintf(ErrorMessageForMethodNotAllowed, method)).Status(http.StatusMethodNotAllowed)
	}))
}
