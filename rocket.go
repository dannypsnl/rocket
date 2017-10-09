package rocket

import (
	"net/http"
	"rocket/routes"
)

type Rocket struct {
	port string
}

func (r *Rocket) Mount(route string, h routes.Handler) *Rocket {
	http.HandleFunc(route+h.Route, h.Handle)
	return r
}

func (r *Rocket) MountNative(route string, handle func(http.ResponseWriter, *http.Request)) *Rocket {
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
