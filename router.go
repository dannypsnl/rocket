package rocket

import (
	"reflect"
	"strings"

	"github.com/dannypsnl/rocket/response"
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
				matchRoute.PathRouteHandler = make(map[string]*handler)
			}
			if _, ok := matchRoute.PathRouteHandler[h.method]; !ok {
				matchRoute.OwnHandler = true
				matchRoute.PathRouteHandler[h.method] = h
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

	matchRoute.OwnHandler = true
	matchRoute.Handlers[h.method] = h
}

func (route *Route) getHandler(requestUrl []string, method string) *handler {
	if len(requestUrl) == 0 {
		if !route.OwnHandler {
			return nil
		}
		if h, ok := route.Handlers[method]; ok {
			return h
		} else {
			return &handler{
				do: reflect.ValueOf(func() *response.Response {
					return response.New("").Status(403)
				}),
			}
		}
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
			if h, ok := next.PathRouteHandler[method]; ok {
				h.addMatchedPathValueIntoContext(requestUrl[i:]...)
				return h
			} else {
				return &handler{
					do: reflect.ValueOf(func() *response.Response {
						return response.New("").Status(403)
					}),
				}
			}
		} else {
			return nil
		}
	}
	if !next.OwnHandler {
		return nil
	}
	if h, ok := next.Handlers[method]; ok {
		return h
	} else {
		return &handler{
			do: reflect.ValueOf(func() *response.Response {
				return response.New("").Status(403)
			}),
		}
	}
}

func isParameter(route string) bool {
	return route[0] == ':'
}
