package main

import (
	"fmt"
	"rocket"
	"rocket/routes"
)

var hello = routes.Handler{
	Route: "/:name/:age",
	Do: func(Context map[string]string) string {
		return fmt.Sprintf("Hello, %s.\nYour age is %s", Context["name"], Context["age"])
	},
}

var src = routes.Handler{
	Route: "/*filepath",
	Do: func(Context map[string]string) string {
		return fmt.Sprintf("filepath: %s", Context["filepath"])
	},
}

func main() {
	fmt.Println("GO web rocket!!!")
	rocket.
		Ignite(":8080").
		Mount("/hello", hello).
		Mount("/src", src).
		Launch()
}
