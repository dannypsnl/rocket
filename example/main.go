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

func homePage() response.Html {
	return `<h1>Title</h1>`
}
func hello(user *User) string {
	return fmt.Sprintf("Hello %s, your age is %d\n", user.Name, user.Age)
}

func main() {
	fmt.Println("GO web rocket!!!")
	rk.
		Ignite(":8080").
		Mount(
			rk.Get("/", homePage),
			rk.Get("/hello/:name/:age", hello),
		).
		Launch()
}
