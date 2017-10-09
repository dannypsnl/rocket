package rocket

import (
	"net/http"
)

type Rocket struct {
	port string
}

func (r *Rocket) Mount(route string, handle func(http.ResponseWriter, *http.Request)) *Rocket {
	http.HandleFunc(route, handle)
	return r
}

func (r *Rocket) Launch() {
	http.ListenAndServe(r.port, nil)
}

func Ignite(port string) *Rocket {
	return &Rocket{
		port: port,
	}
}
