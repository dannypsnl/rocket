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
	open := false
	start := -1
	s := ""
	// TODO: 驗證url之後再綁定，因為url可能含有參數
	// '/:id' is params in url.
	// '/*filepath' is params about filepath.
	// '/home, data' is params from post method.
	for i, r := range route {
		if r == ':' {
			open = true
			start = i + 1
		}
		if open {
			fmt.Println("%s", string(r))
			if r == '/' || i == len(route)-1 {
				s = route[start : i+1]
				open = false
			}
		}
	}
	fmt.Println(s)
	r.handlers[route+h.Route] = h
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
