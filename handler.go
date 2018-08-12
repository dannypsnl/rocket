package rocket

import (
	"reflect"
	"strings"
)

type handler struct {
	route  string
	params map[int]int   // Never custom it. It only for rocket inside.
	do     reflect.Value // do should return response for HTTP writer
	method string
}

func handlerByMethod(route *string, do interface{}, method string) *handler {
	handlerDo := reflect.ValueOf(do)
	h := &handler{
		route:  *route,
		do:     handlerDo,
		method: method,
		params: make(map[int]int),
	}

	handlerT := reflect.TypeOf(do)
	userDefinedT := handlerT.In(0).Elem()
	for idx, r := range strings.Split(h.route, "/")[1:] {
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
