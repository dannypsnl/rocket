package main

import (
	"fmt"
	rk "github.com/dannypsnl/rocket"
)

var hello = rk.Get("/:name/age/:age", func(ctx rk.Ctx) rk.Res {
	return rk.Res(fmt.Sprintf("Hello %s, your age is %s\n", ctx["name"], ctx["age"]))
})

var getDate = rk.Get("/date/:date", func(ctx rk.Ctx) rk.Res {
	return rk.Res(fmt.Sprintf("Time is %s\n", ctx["date"]))
})

var getDateCopy = rk.Get("/date/20", func(ctx rk.Ctx) rk.Res {
	return rk.Res(fmt.Sprintf("Time is 20. Yeee!!!!!!!!!\n"))
})

var index = rk.Get("/", func(ctx rk.Ctx) rk.Res {
	return "index"
})

var static = rk.Get("/*path", func(ctx rk.Ctx) rk.Res {
	return "static"
})

var API = rk.Post("/", func(ctx rk.Ctx) rk.Res {
	return "API"
})

func main() {
	fmt.Println("GO web rocket!!!")
	rk.
		Ignite(":8080").
		Mount("/", index).
		Mount("/", static).
		Mount("/", getDateCopy).
		Mount("/", getDate).
		Mount("/hello", hello).
		Mount("/api", API).
		Launch()
}
