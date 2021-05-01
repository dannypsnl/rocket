package rocket

import (
	"fmt"
	"log"
	"net/http"
	"reflect"

	"github.com/dannypsnl/rocket/response"
	"github.com/dannypsnl/rocket/router"
)

// Rocket is our service.
type Rocket struct {
	port          int
	router        *router.Route
	listOfFairing []Fairing

	allowTLS bool
	// TLS
	certFile string
	keyFile  string
	// cache layer
	defaultHandler reflect.Value
	defaultResp    *response.Response

	// MultiFormBodySizeLimit decide the multiple forms value size
	MultiFormBodySizeLimit int64

	// on close invoke
	onClose func() error
}

// Ignite initial service by port.
func Ignite(port int) *Rocket {
	return &Rocket{
		port: port,
		router: router.New(
			&optionsHandler{},
			createNotAllowHandler,
		),
		listOfFairing: make([]Fairing, 0),
		allowTLS:      false,
		defaultHandler: reflect.ValueOf(func() string {
			return "page not found"
		}),
		// default limit: 10MB
		MultiFormBodySizeLimit: 10,
	}
}

// Attach add fairing to lifecycle for each request and response
func (rk *Rocket) Attach(f Fairing) *Rocket {
	rk.listOfFairing = append(rk.listOfFairing, f)
	return rk
}

// EnableHTTPs would get certFile and keyFile to enable HTTPs
func (rk *Rocket) EnableHTTPs(certFile, keyFile string) *Rocket {
	rk.certFile = certFile
	rk.keyFile = keyFile
	rk.allowTLS = true
	return rk
}

// Mount add handlers into our service.
func (rk *Rocket) Mount(handlers ...*handler) *Rocket {
	for _, h := range handlers {
		rk.router.AddHandler(h.method, h.getRoute(), h)
	}
	return rk
}

// Default receive a function that have signature `func() <T>` for custom response when no route matched,
// <T> means a legal response Type of rocket, e.g. `*response.Response`, `response.Json`
// by default that(status code 404) would returns `"page not found"` when no set this function,
func (rk *Rocket) Default(do interface{}) *Rocket {
	rk.defaultHandler = reflect.ValueOf(do)
	return rk
}

// OnClose takes a function f and runs it after server closed
func (rk *Rocket) OnClose(f func() error) *Rocket {
	rk.onClose = f
	return rk
}

// Launch shoot our service.(start server)
func (rk *Rocket) Launch() {
	for _, f := range rk.listOfFairing {
		f.OnLaunch(rk)
	}
	http.HandleFunc("/", rk.ServeHTTP)
	server := &http.Server{Addr: fmt.Sprintf(":%d", rk.port), Handler: rk}
	defer func() {
		if err := server.Close(); err != nil {
			log.Fatal(err)
		}
		if err := rk.onClose(); err != nil {
			log.Fatal(err)
		}
	}()
	switch {
	case rk.allowTLS:
		log.Fatal(server.ListenAndServeTLS(rk.certFile, rk.keyFile))
	default:
		log.Fatal(server.ListenAndServe())
	}
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
		h.rocket = rk
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
