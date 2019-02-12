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

	guards       []*context.UserContext
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
	return newHandler(reflect.ValueOf(func() *response.Response {
		return response.New(content).Status(code)
	}))
}

func (h *handler) handle(reqURL []string, r *http.Request) *response.Response {
	ctx, err := h.getContexts(reqURL, r)
	if err != nil {
		return response.New(err.Error()).
			Status(http.StatusBadRequest)
	}

	if err := h.verify(reqURL, r); err != nil {
		if err, ok := err.(*VerifyError); ok {
			return response.New(err.Error()).
				Status(err.Status())
		}
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

func (h *handler) Guard(guard Guard) *handler {
	if h.guards == nil {
		h.guards = make([]*context.UserContext, 0)
	}
	contextT := reflect.TypeOf(guard).Elem()
	h.guards = append(h.guards,
		context.
			NewUserContext().
			CacheParamsOffset(contextT, h.routes),
	)
	return h
}

func (h *handler) verify(reqURL []string, r *http.Request) error {
	// no guards
	if h.guards == nil {
		return nil
	}
	ctx, err := h.getGuards(reqURL, r)
	if err != nil {
		return err
	}
	for _, c := range ctx {
		// without check since we already do the static checking by method signature
		guard := c.Interface().(Guard)
		if err := guard.VerifyRequest(); err != nil {
			return err
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

func (h *handler) getContexts(reqURL []string, req *http.Request) ([]reflect.Value, error) {
	return h.fillByCachedUserContexts(h.userContexts, reqURL, req)
}

func (h *handler) getGuards(reqURL []string, req *http.Request) ([]reflect.Value, error) {
	return h.fillByCachedUserContexts(h.guards, reqURL, req)
}

func (h *handler) fillByCachedUserContexts(contexts []*context.UserContext, reqURL []string, req *http.Request) ([]reflect.Value, error) {
	userContexts := make([]reflect.Value, len(contexts))

	req.ParseForm()
	for i, userContext := range contexts {
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
