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

func (rk *Rocket) handler(w http.ResponseWriter, r *http.Request) {
	h := rk.handlers[r.URL.Path]
	fmt.Fprintf(w, h.Do())
}

func (r *Rocket) Mount(route string, h routes.Handler) *Rocket {
	// TODO: 驗證url之後再綁定，因為url可能含有參數
	r.handlers[route+h.Route] = h
	return r
}

func (r *Rocket) Launch() {
	http.HandleFunc("/", r.handler)
	log.Fatal(http.ListenAndServe(r.port, nil))
}

func Ignite(port string) *Rocket {
	return &Rocket{
		port:     port,
		handlers: make(map[string]routes.Handler),
	}
}
