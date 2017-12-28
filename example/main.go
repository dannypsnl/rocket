package main

import (
	"fmt"
	rk "github.com/dannypsnl/rocket"
)

var hello = rk.Get("/:name/age/:age", func(ctx rk.Context) rk.Res {
	return rk.Res(fmt.Sprintf("Hello %s, your age is %s\n", ctx["name"], ctx["age"]))
})

var index = rk.Get("/", func(ctx rk.Context) rk.Res {
	return "index"
})

var static = rk.Get("/*path", func(ctx rk.Context) rk.Res {
	return "static"
})

var API = rk.Post("/", func(ctx rk.Context) rk.Res {
	return "API"
})

func main() {
	fmt.Println("GO web rocket!!!")
	rk.
		Ignite(":8080").
		Mount("/", index).
		Mount("/", static).
		Mount("/hello", hello).
		Mount("/api", API).
		Launch()
}
