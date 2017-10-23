package rocket

type Handler struct {
	Route  string
	params []string // Never custom it. It only for rocket inside.
	Do     func(Context) Response
}

func Get(route string, do func(Context) Response) *Handler {
	return &Handler{
		Route: route,
		Do:    do,
	}
}
