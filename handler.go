package rocket

import "reflect"

type handler struct {
	route  string
	params map[string]string // Never custom it. It only for rocket inside.
	do     reflect.Value     // do should return response for HTTP writer
	method string
}

func handlerByMethod(route *string, do interface{}, method string) *handler {
	handlerT := reflect.TypeOf(do)
	userDefinedT := handlerT.In(0).Elem()
	for i := 0; i < userDefinedT.NumField(); i++ {
		userDefinedT.Field(i).Tag.Get("route")
	}

	handlerDo := reflect.ValueOf(do)
	return &handler{
		route:  *route,
		do:     handlerDo,
		method: method,
		params: make(map[string]string),
	}
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
