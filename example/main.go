package main

import (
	"fmt"
	"github.com/dannypsnl/rocket"
	"io/ioutil"
	"log"
)

var hello = rocket.Handler{
	Route: "/:name/:age",
	Do: func(Context map[string]string) string {
		return fmt.Sprintf("Hello, %s.\nYour age is %s", Context["name"], Context["age"])
	},
}

var index = rocket.Handler{
	Route: "",
	Do: func(Context map[string]string) string {
		path := "./static/index.html"
		buf, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatal(err)
		}
		return string(buf)
	},
}

var src = rocket.Handler{
	Route: "/*filepath",
	Do: func(Context map[string]string) string {
		path := "./static/" + Context["filepath"]
		buf, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatal(err)
		}
		return fmt.Sprintf("%s", string(buf))
	},
}

func main() {
	fmt.Println("GO web rocket!!!")
	rocket.
		Ignite(":8080").
		Mount("/hello", hello).
		Mount("/", index).
		Mount("/", src).
		Launch()
}
