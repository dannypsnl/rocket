package rocket

import (
	"log"
	"net/http"
	"reflect"

	"github.com/dannypsnl/rocket/response"
	"github.com/dannypsnl/rocket/router"
)

// Rocket is our service.
type Rocket struct {
	port          string
	router        *router.Route
	listOfFairing []fairingInterface

	allowTLS bool
	// TLS
	certFile string
	keyFile  string
	// cache layer
	defaultHandler reflect.Value
	defaultResp    *response.Response
}

// Mount add handler into our service.
func (rk *Rocket) Mount(h *handler, hs ...*handler) *Rocket {
	err := rk.router.AddHandler(h.method, h.getRoute(), h)
	if err != nil {
		panic(err)
	}
	for _, h := range hs {
		err := rk.router.AddHandler(h.method, h.getRoute(), h)
		if err != nil {
			panic(err)
		}
	}
	return rk
}

// Attach add fairing to lifecycle for each request and response
func (rk *Rocket) Attach(f fairingInterface) *Rocket {
	rk.listOfFairing = append(rk.listOfFairing, f)
	return rk
}

// Default receive a function that have signature `func() <T>` for custom response when no route matched,
// <T> means a legal response Type of rocket, e.g. `*response.Response`, `response.Json`
// by default that(status code 404) would returns `"page not found"` when no set this function,
func (rk *Rocket) Default(do interface{}) *Rocket {
	rk.defaultHandler = reflect.ValueOf(do)
	return rk
}

// Launch shoot our service.(start server)
func (rk *Rocket) Launch() {
	for _, f := range rk.listOfFairing {
		f.OnLaunch(rk)
	}
	http.HandleFunc("/", rk.ServeHTTP)
	server := &http.Server{Addr: rk.port, Handler: rk}
	defer server.Close()
	switch {
	case rk.allowTLS:
		log.Fatal(server.ListenAndServeTLS(rk.certFile, rk.keyFile))
	default:
		log.Fatal(server.ListenAndServe())
	}
}

// Ignite initial service by port. The format following native HTTP library `:port_number`
func Ignite(port string) *Rocket {
	return &Rocket{
		port: port,
		router: router.New(
			&optionsHandler{},
			createNotAllowHandler,
		),
		listOfFairing: make([]fairingInterface, 0),
		allowTLS:      false,
		defaultHandler: reflect.ValueOf(func() string {
			return "page not found"
		}),
	}
}

// EnableHTTPs would get certFile and keyFile to enable HTTPs
func (rk *Rocket) EnableHTTPs(certFile, keyFile string) *Rocket {
	rk.certFile = certFile
	rk.keyFile = keyFile
	rk.allowTLS = true
	return rk
}

// ServeHTTP is prepare for http server trait, so that you could use `*rocket.Rocket` as `http.Handler`
func (rk *Rocket) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqURL := router.SplitBySlash(r.URL.Path)

	for _, f := range rk.listOfFairing {
		r = f.OnRequest(r)
	}

	// get response
	hand := rk.router.GetHandler(reqURL, r.Method)
	var resp *response.Response
	if h, ok := hand.(*handler); ok && h != nil {
		resp = h.handle(reqURL, r)
	} else {
		resp = rk.defaultResponse()
	}

	for _, f := range rk.listOfFairing {
		resp = f.OnResponse(resp)
	}
	resp.ServeHTTP(w, r)
}

func (rk *Rocket) defaultResponse() *response.Response {
	if rk.defaultResp != nil {
		return rk.defaultResp
	}
	rk.defaultResp = response.New(
		rk.defaultHandler.Call([]reflect.Value{})[0],
	).Status(http.StatusNotFound)
	return rk.defaultResp
}
