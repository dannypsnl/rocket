package main

import (
	"fmt"
	"github.com/dannypsnl/rocket"
)

var hello = rocket.Handler{
	Route: "/:name/:age",
	Do: func(Context map[string]string) string {
		return fmt.Sprintf("Hello, %s.\nYour age is %s\n", "danny", "20") +
			fmt.Sprintf("Hello, %s.\nYour age is %s", Context["name"], Context["age"])
	},
}

func main() {
	fmt.Println("GO web rocket!!!")
	rocket.
		Ignite(":8080").
		Mount("/hello", hello).
		Launch()
}
