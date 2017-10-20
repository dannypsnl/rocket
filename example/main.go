package main

import (
	"fmt"
	"github.com/dannypsnl/rocket"
)

var hello = rocket.Handler{
	Route: "/:name/age/:age",
	Do: func(context rocket.Context) string {
		return fmt.Sprintf("Hello, %s.\nYour age is %s\n", context["name"], context["age"])
	},
}

var index = rocket.Handler{
	Do: func(context rocket.Context) string {
		return "Home"
	},
}

var static = rocket.Handler{
	Route: "/*path",
	Do: func(context rocket.Context) string {
		return "static/" + context["path"]
	},
}

func main() {
	fmt.Println("GO web rocket!!!")
	rocket.
		Ignite(":8080").
		Mount("/", index).
		Mount("/", static).
		Mount("/hello", hello).
		Launch()
}
