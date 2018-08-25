package rocket

import (
	"net/http"
	"reflect"
	"strings"
)

type handler struct {
	routes     []string
	params     map[int]int // Never custom it. It only for rocket inside.
	postParams map[string]int
	do         reflect.Value // do should return response for HTTP writer
	method     string
}

func (h *handler) context(rs []string, req *http.Request) []reflect.Value {
	param := make([]reflect.Value, 0)
	if h.do.Type().NumIn() > 0 {
		contextType := h.do.Type().In(0)
		context := reflect.New(contextType.Elem())

		for idx, route := range h.routes {
			if isParameter(route) {
				param := rs[len(rs)-len(h.routes)+idx]
				index := h.params[idx]
				value := parseParameter(context.Elem().Field(index), param)
				context.Elem().Field(index).
					Set(value)
			}
		}

		req.ParseForm()
		for k, idx := range h.postParams {
			p := req.FormValue(k)
			value := parseParameter(context.Elem().Field(idx), p)
			context.Elem().Field(idx).
				Set(value)
		}

		param = append(param, context)
	}
	return param
}

func handlerByMethod(route *string, do interface{}, method string) *handler {
	handlerDo := reflect.ValueOf(do)
	splitPostParam := strings.Split(*route, ",")
	h := &handler{
		routes:     strings.Split(splitPostParam[0], "/")[1:],
		do:         handlerDo,
		method:     method,
		params:     make(map[int]int),
		postParams: make(map[string]int),
	}

	handlerT := reflect.TypeOf(do)
	if handlerT.NumIn() > 0 {
		userDefinedT := handlerT.In(0).Elem()
		for idx, r := range h.routes {
			if r[0] == ':' {
				for i := 0; i < userDefinedT.NumField(); i++ {
					key := userDefinedT.Field(i).Tag.Get("route")
					if key == r[1:] {
						h.params[idx] = i
						break
					}
				}
			}
		}

		for _, postP := range splitPostParam[1:] {
			for i := 0; i < userDefinedT.NumField(); i++ {
				key := userDefinedT.Field(i).Tag.Get("form")
				if key == postP {
					h.postParams[postP] = i
				}
			}
		}
	}

	return h
}

// Get return a get handler.
func Get(route string, do interface{}) *handler {
	return handlerByMethod(&route, do, "GET")
}

// Post return a post handler.
func Post(route string, do interface{}) *handler {
	return handlerByMethod(&route, do, "POST")
}

// Put return a put handler.
func Put(route string, do interface{}) *handler {
	return handlerByMethod(&route, do, "PUT")
}

// Delete return delete handler.
func Delete(route string, do interface{}) *handler {
	return handlerByMethod(&route, do, "DELETE")
}
