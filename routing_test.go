package rocket_test

import (
	"net/http/httptest"
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
			name:   "root",
			routes: []string{"/"},
			requestRouteMapToMatchedRoute: map[string]string{
				"/": "/",
			},
		},
		{
			name: "multi spec route",
			routes: []string{
				"/a",
				"/b",
				"/a/b",
			},
			requestRouteMapToMatchedRoute: map[string]string{
				"/a":   "/a",
				"/b":   "/b",
				"/a/b": "/a/b",
			},
		},
		{
			name: "variant vs spec",
			routes: []string{
				"/a",
				"/:name",
			},
			requestRouteMapToMatchedRoute: map[string]string{
				"/a": "/a",
				"/b": "/:name",
				"/c": "/:name",
			},
		},
		{
			name: "wildcard vs spec",
			routes: []string{
				"/a",
				"/*wildcard",
			},
			requestRouteMapToMatchedRoute: map[string]string{
				"/a":   "/a",
				"/b":   "/*wildcard",
				"/b/c": "/*wildcard",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			rk := rocket.Ignite("")
			for _, route := range testCase.routes {
				rk.Mount("/", rocket.Get(route,
					func(response string) func() string {
						return func() string { return response }
					}(route)))
			}
			ts := httptest.NewServer(rk)
			defer ts.Close()

			e := httpexpect.New(t, ts.URL)
			for requestRoute, matchedRoute := range testCase.requestRouteMapToMatchedRoute {
				e.GET(requestRoute).Expect().Body().Equal(matchedRoute)
			}
		})
	}
}
