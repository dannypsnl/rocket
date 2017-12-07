package rocket

type handler struct {
	route  string
	params []string // Never custom it. It only for rocket inside.
	do     func(Context) Response
	method string
}

// Get return a get handler.
func Get(route string, do func(Context) Response) *handler {
	return &handler{
		route:  route,
		do:     do,
		method: "GET",
	}
}

// Post return a post handler.
func Post(route string, do func(Context) Response) *handler {
	return &handler{
		route:  route,
		do:     do,
		method: "POST",
	}
}

// Put return a put handler.
func Put(route string, do func(Context) Response) *handler {
	return &handler{
		route:  route,
		do:     do,
		method: "PUT",
	}
}

// Delete return delete handler.
func Delete(route string, do func(Context) Response) *handler {
	return &handler{
		route:  route,
		do:     do,
		method: "DELETE",
	}
}
