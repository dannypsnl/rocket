package main

import (
	"fmt"

	rk "github.com/dannypsnl/rocket"
	"github.com/dannypsnl/rocket/response"
)

type User struct {
	Name string `route:"name"`
	Age  int    `route:"age"`
}

var (
	hello = rk.Get("/:name/age/:age", func(user *User) string {
		return fmt.Sprintf("Hello %s, your age is %d\n", user.Name, user.Age)
	})
	homePage = rk.Get("/", func() response.Html {
		return `<h1>Title</h1>`
	})
)

func main() {
	fmt.Println("GO web rocket!!!")
	rk.
		Ignite(":8080").
		Mount("/", homePage).
		Mount("/hello", hello).
		Launch()
}
