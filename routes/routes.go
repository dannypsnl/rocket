package routes

type Handler struct {
	Route  string
	Params []string // Never custom it. It only for rocket inside.
	Do     func(...interface{}) string
}
