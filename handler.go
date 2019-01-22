package rocket

import (
	"bytes"
	"net/http"
	"reflect"

	"github.com/dannypsnl/rocket/response"
)

type handler struct {
	routes []string
	do     reflect.Value // do should return response for HTTP writer
	method string

	userContextsOffset []int
	cookiesOffset      int
	headerOffset       int

	routeParams       map[int]map[int]int // Never custom it. It only for rocket inside.
	formParams        map[int]map[string]int
	queryParams       map[int]map[string]int
	expectJsonRequest bool

	matchedPath      string
	matchedPathIndex int
}

func newHandler(do reflect.Value) *handler {
	return &handler{
		do:                 do,
		userContextsOffset: make([]int, 0),
		cookiesOffset:      -1,
		headerOffset:       -1,
		matchedPathIndex:   -1,
		routeParams:        make(map[int]map[int]int),
		formParams:         make(map[int]map[string]int),
		queryParams:        make(map[int]map[string]int),
	}
}

func newErrorHandler(code int, content string) *handler {
	h := newHandler(reflect.ValueOf(func() *response.Response {
		return response.New(content).Status(code)
	}))
	return h
}

func (h *handler) Handle(reqURL []string, r *http.Request) *response.Response {
	ctx, err := h.getUserContexts(reqURL, r)
	if err != nil {
		return response.New(err.Error()).
			Status(http.StatusBadRequest)
	}
	resp := h.do.Call(
		ctx,
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

func (h *handler) needCookies() bool {
	return h.cookiesOffset != -1
}
func (h *handler) needHeader() bool {
	return h.headerOffset != -1
}

func (h *handler) getUserContexts(reqURL []string, req *http.Request) ([]reflect.Value, error) {
	userContexts := make([]reflect.Value, h.do.Type().NumIn())

	req.ParseForm()
	for i, offset := range h.userContextsOffset {
		contextT := h.do.Type().In(offset).Elem()
		context := reflect.New(contextT)

		req.ParseForm() // required! Unless we won't get parsed req.Form
		chain := newChain(context).
			pipe(newRouteFiller(
				h.routes,
				reqURL,
				h.routeParams[i],
				h.matchedPathIndex,
				h.matchedPath,
			)).
			pipe(newQueryFiller(h.queryParams[i], req.URL.Query()))
		if h.expectJsonRequest {
			chain.
				pipe(newJSONFiller(req.Body))
		} else {
			chain.
				pipe(newFormFiller(h.formParams[i], req.Form))
		}
		if chain.error() != nil {
			return nil, chain.error()
		}
		userContexts[offset] = context
	}

	if h.needCookies() {
		userContexts[h.cookiesOffset] = reflect.ValueOf(&Cookies{req: req})
	}

	if h.needHeader() {
		userContexts[h.headerOffset] = reflect.ValueOf(&Headers{header: req.Header})
	}

	return userContexts, nil
}
