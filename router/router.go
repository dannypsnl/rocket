package router

import (
	"errors"
	"strings"
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
	wildcardMethodHandlers map[string]Handler
	// methodHandlers stores map Method to handler of this route
	methodHandlers map[string]Handler
	// optionsHandlerBuilder stores a special handler for OPTION method handling
	optionsHandlerBuilder *optionsHandlerBuilder
	optionHandler         OptionsHandler
	notAllowHandler       func(method string) Handler
}

// OptionsHandler interface define an interface that helps router automatically create an OPTION method handler for you
type OptionsHandler interface {
	Build(allowMethods string) Handler
}

type Handler interface {
	WildcardIndex(int)
}

func New(optionsHandler OptionsHandler, notAllowHandler func(method string) Handler) *Route {
	return &Route{
		children:        make(map[string]*Route),
		methodHandlers:  make(map[string]Handler),
		optionHandler:   optionsHandler,
		notAllowHandler: notAllowHandler,
	}
}

func (r *Route) mustGetVariableRoute() *Route {
	if r.variableRoute == nil {
		r.variableRoute = New(r.optionHandler, r.notAllowHandler)
	}
	return r.variableRoute
}
func (r *Route) mustGetWildcardRoute() map[string]Handler {
	if r.wildcardMethodHandlers == nil {
		r.wildcardMethodHandlers = make(map[string]Handler)
	}
	return r.wildcardMethodHandlers
}

var PanicDuplicateRoute = errors.New("duplicate route")

func SplitBySlash(routeStr string) []string {
	route := make([]string, 0)
	for _, r := range strings.Split(strings.Trim(routeStr, "/"), "/") {
		if r != "" {
			route = append(route, r)
		}
	}
	return route
}

// AddHandler bind method+route with handler, notice that it would just panic when get duplicate route
func (root *Route) AddHandler(method, route string, h Handler) {
	fullRoute := SplitBySlash(route)
	curRoute := root
	for i, r := range fullRoute {
		switch r[0] {
		case ':':
			curRoute = curRoute.mustGetVariableRoute()
			continue
		case '*':
			h.WildcardIndex(i)
			wildcard := curRoute.mustGetWildcardRoute()
			if _, routeExisted := wildcard[method]; routeExisted {
				panic(PanicDuplicateRoute)
			}
			curRoute.addHandlerOn(method, wildcard, h)
			return
		default:
			if _, ok := curRoute.children[r]; !ok {
				curRoute.children[r] = New(root.optionHandler, root.notAllowHandler)
			}
		}
		curRoute = curRoute.children[r]
	}

	if _, sameRouteExisted := curRoute.methodHandlers[method]; sameRouteExisted {
		panic(PanicDuplicateRoute)
	}
	curRoute.addHandlerOn(method, curRoute.methodHandlers, h)
}

func (r *Route) addHandlerOn(method string, m map[string]Handler, h Handler) {
	if r.optionsHandlerBuilder == nil {
		r.optionsHandlerBuilder = newOptionsHandler()
	}
	r.optionsHandlerBuilder.addMethod(method)
	m["OPTIONS"] = r.optionHandler.Build(r.optionsHandlerBuilder.build())
	m[method] = h
	r.ownHandler = true
}

func (r *Route) GetHandler(requestUrl []string, method string) Handler {
	if len(requestUrl) == 0 {
		if !r.ownHandler {
			return nil
		}
		if h, ok := r.methodHandlers[method]; ok {
			return h
		}
		return r.notAllowHandler(method)
	}

	headOfUrl, onRestUrl := requestUrl[0], requestUrl[1:]
	if router, ok := r.children[headOfUrl]; ok {
		if h := router.GetHandler(onRestUrl, method); h != nil {
			return h
		}
	}
	if r.variableRoute != nil {
		if h := r.variableRoute.GetHandler(onRestUrl, method); h != nil {
			return h
		}
	}
	if r.wildcardMethodHandlers != nil {
		if h, hasWildcardRouteHandler := r.wildcardMethodHandlers[method]; hasWildcardRouteHandler {
			return h
		}
		return r.notAllowHandler(method)
	} else {
		return nil
	}
}

type optionsHandlerBuilder struct {
	methods []string
}

func newOptionsHandler() *optionsHandlerBuilder {
	return &optionsHandlerBuilder{
		methods: make([]string, 0),
	}
}
func (o *optionsHandlerBuilder) addMethod(method string) {
	o.methods = append(o.methods, method)
}
func (o *optionsHandlerBuilder) build() string {
	allowMethods := "OPTIONS"
	for _, m := range o.methods {
		allowMethods += ", " + m
	}
	return allowMethods
}
