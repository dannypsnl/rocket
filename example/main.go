package main

import (
	"fmt"
	"log"
	"os"
	"rocket"
)

var hello = rocket.Handler{
	Route: "/:name/:age",
	Do: func(Context map[string]string) string {
		return fmt.Sprintf("Hello, %s.\nYour age is %s", Context["name"], Context["age"])
	},
}

var src = rocket.Handler{
	Route: "/*filepath",
	Do: func(Context map[string]string) string {
		f, err := os.OpenFile("./static"+Context["filepath"], os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			log.Fatal(err)
		}
		var bs []byte
		var file string
		i, err := f.Read(bs)
		for i != 0 {
			file += string(bs)
			i, err = f.Read(bs)
		}
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
		return fmt.Sprintf("%s", file)
	},
}

func main() {
	fmt.Println("GO web rocket!!!")
	rocket.
		Ignite(":8080").
		Mount("/hello", hello).
		Mount("/", src).
		Launch()
}
