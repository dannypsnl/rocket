package rocket

type handler struct {
	route  string
	params []string // Never custom it. It only for rocket inside.
	do     func(Context) Response
	method string
}

func Get(route string, do func(Context) Response) *handler {
	return &handler{
		route:  route,
		do:     do,
		method: "GET",
	}
}

func Post(route string, do func(Context) Response) *handler {
	return &handler{
		route:  route,
		do:     do,
		method: "POST",
	}
}
