package rocket

import (
	"fmt"
	"log"
	"net/http"

	"rocket/routes"
)

type Rocket struct {
	port     string
	handlers map[string]routes.Handler
}

func (r *Rocket) Mount(route string, h routes.Handler) *Rocket {
	route += h.Route
	match := ""

	firstTime := true
	start := 0
	s := ""
	// TODO: 驗證url之後再綁定，因為url可能含有參數
	// '/:id' is params in url.
	// '/*filepath' is params about filepath.
	// '/home, data' is params from post method.
	for i, r := range route {
		if r == ':' || r == '*' {
			if firstTime {
				match = route[:i-1]
				firstTime = false
			}
			start = i + 1
		}
		if i == len(route)-1 {
			s = route[start : i+1]
			fmt.Println(s)
		}
		if r == '/' {
			// Get param setting string.
			s = route[start:i]
			fmt.Println(s)
		}
	}
	r.handlers[match] = h
	return r
}

func (rk *Rocket) Launch() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		h := rk.handlers[r.URL.Path]
		fmt.Fprintf(w, h.Do())
	})
	log.Fatal(http.ListenAndServe(rk.port, nil))
}

func Ignite(port string) *Rocket {
	return &Rocket{
		port:     port,
		handlers: make(map[string]routes.Handler),
	}
}
