package rocket_test

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dannypsnl/rocket"

	"github.com/gavv/httpexpect"
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
				"GET;/a":      "/a",
				"GET;/b":      "/*wildcard",
				"GET;/b/c":    "/*wildcard",
				"GET;/a/b/c":  "/*wildcard",
				"POST;/a/b/c": "request resource does not support http method 'POST'",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			rk := rocket.Ignite("")
			for _, route := range testCase.routes {
				handler := rocket.Get("/", func() string { return "wrong" })
				method, route := splitRuleToMethodRoute(route)
				handleFunc := func(response string) func() string {
					return func() string { return response }
				}(route)

				switch method {
				case "GET":
					handler = rocket.Get(route, handleFunc)
				case "POST":
					handler = rocket.Post(route, handleFunc)
				}
				rk.Mount("/", handler)
			}
			ts := httptest.NewServer(rk)
			defer ts.Close()

			e := httpexpect.New(t, ts.URL)
			for requestRoute, matchedRoute := range testCase.requestRouteMapToMatchedRoute {
				method, route := splitRuleToMethodRoute(requestRoute)
				e.Request(method, route).Expect().Body().Equal(matchedRoute)
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
