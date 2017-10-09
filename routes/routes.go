package routes

type Handler struct {
	Route string
	Do    func(...interface{}) string
}
