package rocket

type Handler struct {
	Route  string
	Params []string // Never custom it. It only for rocket inside.
	Do     func(map[string]string) string
}
