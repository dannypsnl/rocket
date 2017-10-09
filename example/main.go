package main

import (
	"fmt"
	"net/http"
	"rocket"
	"rocket/routes"
)

var hello = routes.Handler{
	Route: "",
	Do: func(...interface{}) string {
		return "Hello"
	},
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello index")
}

func main() {
	fmt.Println("GO web rocket!!!")
	rocket.
		Ignite(":8080").
		Mount("/hello", hello).
		Launch()
}
