package rocket

type handler struct {
	route  string
	params []string // Never custom it. It only for rocket inside.
	do     handleMethod
	method string
}

type handleMethod func(Context) Res

func handlerByMethod(route *string, do handleMethod, method string) *handler {
	return &handler{
		route:  *route,
		do:     do,
		method: method,
	}
}

// Get return a get handler.
func Get(route string, do handleMethod) *handler {
	return handlerByMethod(&route, do, "GET")
}

// Post return a post handler.
func Post(route string, do handleMethod) *handler {
	return handlerByMethod(&route, do, "POST")
}

// Put return a put handler.
func Put(route string, do handleMethod) *handler {
	return handlerByMethod(&route, do, "PUT")
}

// Delete return delete handler.
func Delete(route string, do handleMethod) *handler {
	return handlerByMethod(&route, do, "DELETE")
}
