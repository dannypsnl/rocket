package routes

type Handler struct {
	Route string
	Match []string
	Do    func(...interface{}) string
}
