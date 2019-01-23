package rocket

import (
	"reflect"

	"github.com/dannypsnl/rocket/internal"
)

func handlerByMethod(route *string, do interface{}, method string) *handler {
	handlerDo := reflect.ValueOf(do)
	h := newHandler(handlerDo)
	h.method = method

	h.routes = splitBySlash(*route)

	handlerFuncT := reflect.TypeOf(do)
	h.userContexts = make([]*internal.UserContext, handlerFuncT.NumIn())

	for i := 0; i < handlerFuncT.NumIn(); i++ {
		t := handlerFuncT.In(i).Elem()
		userContext := internal.NewUserContext()
		switch {
		case t.AssignableTo(reflect.TypeOf(Cookies{})):
			userContext.IsCookies = true
		case t.AssignableTo(reflect.TypeOf(Headers{})):
			userContext.IsHeaders = true
		default:
			// We not sure what is it, so just assume it's user defined context
			contextT := handlerFuncT.In(i).Elem()
			userContext.CacheParamsOffset(contextT, h.routes)
		}
		h.userContexts[i] = userContext
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
