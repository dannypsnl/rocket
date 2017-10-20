package main

import (
	"fmt"
	"github.com/dannypsnl/rocket"
)

var hello = rocket.Handler{
	Route: "/:name/age/:age",
	Do: func(ctx rocket.Context) string {
		return fmt.Sprintf("Hello, %s.\nYour age is %s\n", ctx["name"], ctx["age"])
	},
}

var index = rocket.Handler{
	Do: func(ctx rocket.Context) string {
		return "Home"
	},
}

var static = rocket.Handler{
	Route: "/*path",
	Do: func(ctx rocket.Context) string {
		return "static/" + ctx["path"]
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
