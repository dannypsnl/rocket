package main

import (
	"fmt"
	rk "github.com/dannypsnl/rocket"
)

var hello = rk.Get("/:name/age/:age", func(ctx rk.Context) rk.Response {
	return "hello"
})

var index = rk.Get("/", func(ctx rk.Context) rk.Response {
	return "index"
})

var static = rk.Get("/*path", func(ctx rk.Context) rk.Response {
	return "static"
})

func main() {
	fmt.Println("GO web rocket!!!")
	rk.
		Ignite(":8080").
		Mount("/", index).
		Mount("/", static).
		Mount("/hello", hello).
		Launch()
}
