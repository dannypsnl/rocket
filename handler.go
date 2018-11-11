package rocket

import (
	"bytes"
	"errors"
	"net/http"
	"reflect"

	"github.com/dannypsnl/rocket/response"
)

type handler struct {
	routes []string
	do     reflect.Value // do should return response for HTTP writer
	method string

	userDefinedContextOffset int
	cookiesOffset            int
	headerOffset             int

	routeParams       map[int]int // Never custom it. It only for rocket inside.
	formParams        map[string]int
	queryParams       map[string]int
	expectJsonRequest bool

	matchedPath      string
	matchedPathIndex int
}

func newHandler(do reflect.Value) *handler {
	return &handler{
		do: do,
		userDefinedContextOffset: -1,
		cookiesOffset:            -1,
		headerOffset:             -1,
		matchedPathIndex:         -1,
	}
}

func newErrorHandler(code int, content string) *handler {
	h := newHandler(reflect.ValueOf(func() *response.Response {
		return response.New(content).Status(code)
	}))
	return h
}

func (h *handler) Handle(reqURL []string, r *http.Request) *response.Response {
	resp := h.do.Call(
		h.context(reqURL, r),
	)[0].Interface()

	switch v := resp.(type) {
	case *response.Response:
		return v
	default:
		return response.New(v)
	}
}

func (h *handler) addMatchedPathValueIntoContext(paths ...string) {
	path := bytes.NewBuffer([]byte(``))
	for _, v := range paths {
		path.WriteString(v)
		path.WriteRune('/')
	}
	h.matchedPath = path.String()[:path.Len()-1]
}

func (h *handler) hasUserDefinedContext() bool {
	return h.userDefinedContextOffset != -1
}
func (h *handler) needCookies() bool {
	return h.cookiesOffset != -1
}
func (h *handler) needHeader() bool {
	return h.headerOffset != -1
}

func (h *handler) context(reqURL []string, req *http.Request) []reflect.Value {
	param := make([]reflect.Value, h.do.Type().NumIn())
	if h.hasUserDefinedContext() {
		contextType := h.do.Type().In(h.userDefinedContextOffset).Elem()
		context := reflect.New(contextType)

		req.ParseForm() // required! Unless we won't get parsed req.Form
		chain := newChain(context).
			pipe(newRouteFiller(
				h.routes,
				reqURL,
				h.routeParams,
				h.matchedPathIndex,
				h.matchedPath,
			)).
			pipe(newQueryFiller(h.queryParams, req.URL.Query()))
		if h.expectJsonRequest {
			chain.
				pipe(newJSONFiller(req.Body))
		} else {
			chain.
				pipe(newFormFiller(h.formParams, req.Form))
		}
		if chain.error() != nil {
			param[h.userDefinedContextOffset] = reflect.ValueOf(errors.New("400"))
			return param
		}
		param[h.userDefinedContextOffset] = context
	}

	if h.needCookies() {
		param[h.cookiesOffset] = reflect.ValueOf(&Cookies{req: req})
	}

	if h.needHeader() {
		param[h.headerOffset] = reflect.ValueOf(&Headers{header: req.Header})
	}

	return param
}
