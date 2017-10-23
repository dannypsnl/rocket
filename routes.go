package rocket

type Handler struct {
	Route  string
	params []string // Never custom it. It only for rocket inside.
	Do     func(Context) Response
}
