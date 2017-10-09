package routes

import (
	"fmt"
	"net/http"
)

type Handler struct {
	Route string
	Do    func(...interface{}) string
}

func (rh *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	res := rh.Do()
	fmt.Fprintf(w, res)
}
