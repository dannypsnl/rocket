package rocket

import (
	"reflect"
	"strings"
)

func handlerByMethod(route *string, do interface{}, method string) *handler {
	handlerDo := reflect.ValueOf(do)
	h := newHandler(handlerDo)
	h.method = method
	h.routes = strings.Split(strings.Trim(*route, "/"), "/")
	h.routeParams = make(map[int]int)
	h.formParams = make(map[string]int)
	h.queryParams = make(map[string]int)

	handlerFuncT := reflect.TypeOf(do)

	for i := 0; i < handlerFuncT.NumIn(); i++ {
		t := handlerFuncT.In(i).Elem()
		switch {
		case t.AssignableTo(reflect.TypeOf(Cookies{})):
			h.cookiesOffset = i
		case t.AssignableTo(reflect.TypeOf(Headers{})):
			h.headerOffset = i
		default:
			// We not sure what is it, so just assume it's user defined context
			h.userDefinedContextOffset = i
		}
	}

	if h.userDefinedContextOffset != -1 {
		contextT := handlerFuncT.In(h.userDefinedContextOffset).Elem()

		routeParams := make(map[string]int)
		for i := 0; i < contextT.NumField(); i++ {
			tagOfField := contextT.Field(i).Tag
			key, ok := tagOfField.Lookup("route")
			if ok {
				routeParams[key] = i
			}
			key, ok = tagOfField.Lookup("form")
			if ok {
				h.formParams[key] = i
			}
			key, ok = tagOfField.Lookup("query")
			if ok {
				h.queryParams[key] = i
			}
			_, ok = tagOfField.Lookup("json")
			if !h.expectJsonRequest && ok {
				h.expectJsonRequest = ok
			}
		}

		for idx, r := range h.routes {
			// a route part like `:name`
			if r[0] == ':' || r[0] == '*' {
				// r[1:] is `name`, that's the key we expected
				h.routeParams[idx] = routeParams[r[1:]]
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

// Patch return a patch handler.
func Patch(route string, do interface{}) *handler {
	return handlerByMethod(&route, do, "PATCH")
}

// Delete return delete handler.
func Delete(route string, do interface{}) *handler {
	return handlerByMethod(&route, do, "DELETE")
}
