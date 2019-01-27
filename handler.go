package rocket

import (
	"bytes"
	"net/http"
	"reflect"

	"github.com/dannypsnl/rocket/internal/context"
	"github.com/dannypsnl/rocket/response"
)

type handler struct {
	routes []string
	do     reflect.Value // do should return response for HTTP writer
	method string

	userContexts []*context.UserContext

	matchedPath      string
	matchedPathIndex int
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

func (h *handler) handle(reqURL []string, r *http.Request) *response.Response {
	ctx, err := h.getUserContexts(reqURL, r)
	if err != nil {
		return response.New(err.Error()).
			Status(http.StatusBadRequest)
	}

	if err := h.verify(ctx, r); err != nil {
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

func (h *handler) verify(ctx []reflect.Value, r *http.Request) error {
	for _, c := range ctx {
		guard, isGuard := c.Interface().(Guard)
		if isGuard {
			if err := guard.VerifyRequest(); err != nil {
				return err
			}
		}
	}
	return nil
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
		basicChain := []filler{
			newRouteFiller(
				h.routes,
				reqURL,
				userContext.RouteParams,
				h.matchedPathIndex,
				h.matchedPath,
			),
			newQueryFiller(userContext.QueryParams, req.URL.Query()),
		}
		if userContext.ExpectJSONRequest {
			basicChain = append(basicChain, newJSONFiller(req.Body))
		} else {
			basicChain = append(basicChain, newFormFiller(userContext.FormParams, req.Form))
		}
		if userContext.ExpectCookies() {
			basicChain = append(basicChain, newCookiesFiller(userContext.CookiesParams, req))
		}
		if userContext.ExpectHeader() {
			basicChain = append(basicChain, newHeaderFiller(userContext.HeaderParams, req.Header))
		}
		ctx := reflect.New(userContext.ContextType)
		for _, filler := range basicChain {
			if err := filler.fill(ctx); err != nil {
				return nil, err
			}
		}
		userContexts[i] = ctx
	}
	return userContexts, nil
}

type optionsHandler struct {
	methods []string
}

func newOptionsHandler() *optionsHandler {
	return &optionsHandler{
		methods: make([]string, 0),
	}
}

func (o *optionsHandler) addMethod(method string) *optionsHandler {
	o.methods = append(o.methods, method)
	return o
}

func (o *optionsHandler) build() *handler {
	allowMethods := "OPTIONS"
	for _, m := range o.methods {
		allowMethods += ", " + m
	}
	return newHandler(reflect.ValueOf(func() *response.Response {
		return response.New("").
			Headers(response.Headers{
				"Allow": allowMethods,
			})
	}))
}
