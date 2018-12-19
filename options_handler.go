package rocket

import (
	"reflect"

	"github.com/dannypsnl/rocket/response"
)

type optionsHandler struct {
	methods []string
}

func newOptionsHandler() *optionsHandler {
	return &optionsHandler{
		methods: make([]string, 0),
	}
}

func (o *optionsHandler) addMethod(method string) *optionsHandler {
	o.methods = append(o.methods, method)
	return o
}

func (o *optionsHandler) build() *handler {
	allowMethods := "OPTIONS"
	for _, m := range o.methods {
		allowMethods += ", " + m
	}
	return newHandler(reflect.ValueOf(func() *response.Response {
		return response.New("").
			Headers(response.Headers{
				"Allow": allowMethods,
			})
	}))
}
