package rocket

import (
	"fmt"
	"log"
	"net/http"

	"rocket/routes"
)

func (rk *Rocket) handler(w http.ResponseWriter, r *http.Request) {
	h := rk.routes[r.URL.Path]
	fmt.Fprintf(w, h.Do())
}

type Rocket struct {
	port   string
	routes map[string]routes.Handler
}

func (r *Rocket) Mount(route string, h routes.Handler) *Rocket {
	// TODO: 驗證url之後再綁定，因為url可能含有參數
	r.routes[route+h.Route] = h
	return r
}

func (r *Rocket) Launch() {
	http.HandleFunc("/", r.handler)
	log.Fatal(http.ListenAndServe(r.port, nil))
}

func Ignite(port string) *Rocket {
	return &Rocket{
		port:   port,
		routes: make(map[string]routes.Handler),
	}
}
