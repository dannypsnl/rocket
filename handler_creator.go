package rocket

import (
	"reflect"

	"github.com/dannypsnl/rocket/internal/context"
	"github.com/dannypsnl/rocket/router"
)

func handlerByMethod(route *string, do interface{}, method string) *handler {
	handlerDo := reflect.ValueOf(do)
	h := newHandler(handlerDo)
	h.method = method

	h.route = *route
	h.routes = router.SplitBySlash(*route)

	handlerFuncT := reflect.TypeOf(do)
	h.userContexts = make([]*context.UserContext, handlerFuncT.NumIn())

	for i := 0; i < handlerFuncT.NumIn(); i++ {
		contextT := handlerFuncT.In(i).Elem()
		h.userContexts[i] = context.NewUserContext().
			CacheParamsOffset(contextT, h.routes)
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
