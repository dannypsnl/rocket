package main

import (
	"fmt"
	"net/http"
	"rocket"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello index")
}

func main() {
	fmt.Println("GO web rocket!!!")
	rocket.
		Ignite(":8080").
		Mount("/", IndexHandler).
		Launch()
}
