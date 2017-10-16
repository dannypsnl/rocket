package main

import (
	"fmt"
	"github.com/dannypsnl/rocket"
)

var hello = rocket.Handler{
	Route: "/:name/age/:age",
	Do: func(Context map[string]string) string {
		return fmt.Sprintf("Hello, %s.\nYour age is %s\n", Context["name"], Context["age"])
	},
}

var index = rocket.Handler{
	Do: func(Context map[string]string) string {
		return "Home"
	},
}

var static = rocket.Handler{
	Route: "/*path",
	Do: func(Context map[string]string) string {
		return "Home"
	},
}

func main() {
	fmt.Println("GO web rocket!!!")
	rocket.
		Ignite(":8080").
		Mount("/hello", hello).
		Mount("/", static).
		Mount("/", index).
		Launch()
}
