package main

import (
	"fmt"
	"rocket"
	"rocket/routes"
)

var hello = routes.Handler{
	Route: "",
	Do: func(...interface{}) string {
		return "Hello"
	},
}

var src = routes.Handler{
	Route: "",
	Do: func(...interface{}) string {
		return "Hi"
	},
}

func main() {
	fmt.Println("GO web rocket!!!")
	rocket.
		Ignite(":8080").
		Mount("/hello/:name/:age", hello).
		Mount("/hi/danny/:age", hello).
		Mount("/src/*filepath", src).
		Launch()
}
