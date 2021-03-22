package router_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/dannypsnl/rocket/router"

	"github.com/stretchr/testify/require"
)

func TestRouting(t *testing.T) {
	testCases := []struct {
		name                          string
		routes                        []string
		requestRouteMapToMatchedRoute map[string]string
	}{
		{
			name:   "no route",
			routes: make([]string, 0),
			requestRouteMapToMatchedRoute: map[string]string{
				"GET;/":  "page not found",
				"GET;/a": "page not found",
			},
		},
		{
			name:   "root",
			routes: []string{"GET;/"},
			requestRouteMapToMatchedRoute: map[string]string{
				"GET;/": "/",
			},
		},
		{
			name: "multi spec route",
			routes: []string{
				"GET;/a",
				"GET;/b",
				"GET;/a/b",
			},
			requestRouteMapToMatchedRoute: map[string]string{
				"GET;/a":   "/a",
				"GET;/b":   "/b",
				"GET;/a/b": "/a/b",
			},
		},
		{
			name: "variant vs spec",
			routes: []string{
				"GET;/a",
				"GET;/:name",
			},
			requestRouteMapToMatchedRoute: map[string]string{
				"GET;/a": "/a",
				"GET;/b": "/:name",
				"GET;/c": "/:name",
			},
		},
		{
			name: "wildcard vs spec",
			routes: []string{
				"GET;/a",
				"GET;/*wildcard",
			},
			requestRouteMapToMatchedRoute: map[string]string{
				"GET;/a":                    "/a",
				"GET;/b":                    "/*wildcard",
				"GET;/b/c":                  "/*wildcard",
				"GET;/a/b/c/d/e/d/das/cast": "/*wildcard",
				"POST;/a/b/c":               "request resource does not support http method 'POST'",
				"POST;/a":                   "request resource does not support http method 'POST'",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			r := router.New(&optionsHandler{}, notAllowHandler)
			for _, requestRoute := range testCase.routes {
				method, route := splitRuleToMethodRoute(requestRoute)
				h := &handler{
					route: route,
				}
				r.AddHandler(method, route, h)
			}

			for requestRoute, matchedRoute := range testCase.requestRouteMapToMatchedRoute {
				method, route := splitRuleToMethodRoute(requestRoute)
				h := r.GetHandler(router.SplitBySlash(route), method)
				if len(router.SplitBySlash(matchedRoute)) > 0 && router.SplitBySlash(matchedRoute)[0] == matchedRoute {
					if h != nil {
						// method not allow
						require.Equal(t, matchedRoute, h.(*handler).message)
					} else {
						// page not found
						require.Equal(t, nil, h)
					}
				} else if h != nil {
					require.Equal(t, matchedRoute, h.(*handler).route)
				}
			}
		})
	}
}

func splitRuleToMethodRoute(rule string) (method string, route string) {
	ss := strings.Split(rule, ";")
	method = ss[0]
	route = ss[1]
	return
}

func notAllowHandler(method string) router.Handler {
	return &handler{
		message: fmt.Sprintf("request resource does not support http method '%s'", method),
	}
}

type optionsHandler struct{}

func (o *optionsHandler) Build(allowMethods string) router.Handler {
	return nil
}

type handler struct {
	route   string
	message string
}

func (h *handler) WildcardIndex(i int) {}
