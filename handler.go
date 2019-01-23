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

	userContexts []*UserContext

	matchedPath      string
	matchedPathIndex int
}

type UserContext struct {
	contextType       reflect.Type
	isCookies         bool
	isHeaders         bool
	routeParams       map[int]int
	formParams        map[string]int
	queryParams       map[string]int
	expectJSONRequest bool
}

func newUserContext() *UserContext {
	return &UserContext{
		isCookies:         false,
		isHeaders:         false,
		routeParams:       make(map[int]int),
		formParams:        make(map[string]int),
		queryParams:       make(map[string]int),
		expectJSONRequest: false,
	}
}

func (ctx *UserContext) cacheParamsOffset(contextT reflect.Type, routes []string) {
	ctx.contextType = contextT
	routeParams := make(map[string]int)
	for i := 0; i < contextT.NumField(); i++ {
		tagOfField := contextT.Field(i).Tag
		key, ok := tagOfField.Lookup("route")
		if ok {
			routeParams[key] = i
		}
		key, ok = tagOfField.Lookup("form")
		if ok {
			ctx.formParams[key] = i
		}
		key, ok = tagOfField.Lookup("query")
		if ok {
			ctx.queryParams[key] = i
		}
		_, ok = tagOfField.Lookup("json")
		if !ctx.expectJSONRequest && ok {
			ctx.expectJSONRequest = ok
		}
	}

	for idx, r := range routes {
		// a route part like `:name`
		if r[0] == ':' || r[0] == '*' {
			// r[1:] is `name`, that's the key we expected
			param := r[1:]
			if _, ok := routeParams[param]; ok {
				ctx.routeParams[idx] = routeParams[param]
			}
		}
	}
}

func newHandler(do reflect.Value) *handler {
	return &handler{
		do:               do,
		matchedPathIndex: -1,
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

func (h *handler) getUserContexts(reqURL []string, req *http.Request) ([]reflect.Value, error) {
	userContexts := make([]reflect.Value, h.do.Type().NumIn())

	req.ParseForm()
	for i, userContext := range h.userContexts {
		if userContext.isCookies {
			userContexts[i] = reflect.ValueOf(&Cookies{req: req})
		} else if userContext.isHeaders {
			userContexts[i] = reflect.ValueOf(&Headers{header: req.Header})
		} else {
			context := reflect.New(userContext.contextType)
			chain := newChain(context).
				pipe(newRouteFiller(
					h.routes,
					reqURL,
					userContext.routeParams,
					h.matchedPathIndex,
					h.matchedPath,
				)).
				pipe(newQueryFiller(userContext.queryParams, req.URL.Query()))
			if userContext.expectJSONRequest {
				chain.
					pipe(newJSONFiller(req.Body))
			} else {
				chain.
					pipe(newFormFiller(userContext.formParams, req.Form))
			}
			if chain.error() != nil {
				return nil, chain.error()
			}
			userContexts[i] = context
		}
	}
	return userContexts, nil
}
