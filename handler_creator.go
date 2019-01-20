package rocket

import (
	"reflect"
)

func handlerByMethod(route *string, do interface{}, method string) *handler {
	handlerDo := reflect.ValueOf(do)
	h := newHandler(handlerDo)
	h.method = method

	h.routes = splitBySlash(*route)

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
			h.userContextsOffset = append(h.userContextsOffset, i)
		}
	}

	for i, contextOffset := range h.userContextsOffset {
		contextT := handlerFuncT.In(contextOffset).Elem()
		h.routeParams[i] = make(map[int]int)
		h.formParams[i] = make(map[string]int)
		h.queryParams[i] = make(map[string]int)

		routeParams := make(map[string]int)
		for j := 0; j < contextT.NumField(); j++ {
			tagOfField := contextT.Field(j).Tag
			key, ok := tagOfField.Lookup("route")
			if ok {
				routeParams[key] = j
			}
			key, ok = tagOfField.Lookup("form")
			if ok {
				h.formParams[i][key] = j
			}
			key, ok = tagOfField.Lookup("query")
			if ok {
				h.queryParams[i][key] = j
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
				param := r[1:]
				if _, ok := routeParams[param]; ok {
					h.routeParams[i][idx] = routeParams[param]
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

// Patch return a patch handler.
func Patch(route string, do interface{}) *handler {
	return handlerByMethod(&route, do, "PATCH")
}

// Delete return delete handler.
func Delete(route string, do interface{}) *handler {
	return handlerByMethod(&route, do, "DELETE")
}
