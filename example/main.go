package main

import (
	"fmt"
	rk "github.com/dannypsnl/rocket"
)

type User struct {
	Name string `route:"name"`
	Age  int    `route:"age"`
}

var hello = rk.Get("/:name/age/:age", func(user User) rk.Res {
	return rk.Res(fmt.Sprintf("Hello %s, your age is %s\n", user.Name, user.Age))
})

func main() {
	fmt.Println("GO web rocket!!!")
	rk.
		Ignite(":8080").
		Mount("/hello", hello).
		Launch()
}
