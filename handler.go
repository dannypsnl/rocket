package rocket

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/dannypsnl/rocket/internal/context"
	"github.com/dannypsnl/rocket/response"
	"github.com/dannypsnl/rocket/router"
)

// Get return a get handler.
func Get(route string, do interface{}) *handler {
	return handlerByMethod(&route, do, "GET")
}

// Post return a post handler.
func Post(route string, do interface{}) *handler {
	return handlerByMethod(&route, do, "POST")
}

// Put return a put handler.
func Put(route string, do interface{}) *handler {
	return handlerByMethod(&route, do, "PUT")
}

// Patch return a patch handler.
func Patch(route string, do interface{}) *handler {
	return handlerByMethod(&route, do, "PATCH")
}

// Delete return delete handler.
func Delete(route string, do interface{}) *handler {
	return handlerByMethod(&route, do, "DELETE")
}

func handlerByMethod(route *string, do interface{}, method string) *handler {
	handlerDo := reflect.ValueOf(do)
	h := newHandler(handlerDo)
	h.method = method

	h.route = *route
	h.routes = router.SplitBySlash(*route)

	handlerFuncT := reflect.TypeOf(do)
	h.userContexts = make([]*context.UserContext, handlerFuncT.NumIn())

	for i := 0; i < handlerFuncT.NumIn(); i++ {
		contextT := handlerFuncT.In(i).Elem()
		h.userContexts[i] = context.NewUserContext().
			CacheParamsOffset(contextT, h.routes)
	}
	return h
}

type handler struct {
	route  string
	routes []string
	do     reflect.Value // do should return response for HTTP writer
	method string

	guards       []*context.UserContext
	userContexts []*context.UserContext

	wildcardIndex int
}

func newHandler(do reflect.Value) *handler {
	if do.Type().NumOut() <= 0 {
		panic(fmt.Sprintf("handling function should be non-void function but got: %s", do.Type()))
	}
	return &handler{
		do:            do,
		wildcardIndex: -1,
	}
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
			Status(http.StatusInternalServerError)
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

func (h *handler) getContexts(reqURL []string, req *http.Request) ([]reflect.Value, error) {
	return h.fillByCachedUserContexts(h.userContexts, reqURL, req)
}

func (h *handler) getGuards(reqURL []string, req *http.Request) ([]reflect.Value, error) {
	return h.fillByCachedUserContexts(h.guards, reqURL, req)
}

func (h *handler) fillByCachedUserContexts(contexts []*context.UserContext, reqURL []string, req *http.Request) ([]reflect.Value, error) {
	userContexts := make([]reflect.Value, len(contexts))

	err := req.ParseForm()
	if err != nil {
		return nil, err
	}
	for i, userContext := range contexts {
		basicChain := []filler{
			newRouteFiller(
				h.routes,
				reqURL,
				userContext.RouteParams,
				h.wildcardIndex,
			),
			newQueryFiller(userContext.QueryParams, req.URL.Query()),
		}
		if userContext.ExpectJSONRequest {
			basicChain = append(basicChain, newJSONFiller(req.Body))
		} else if userContext.ExpectMultiFormsRequest {
			basicChain = append(basicChain, newMultiFormFiller(userContext.MultiFormParams, userContext.MultiFormParamsIsFile, req))
		} else {
			basicChain = append(basicChain, newFormFiller(userContext.FormParams, req.Form))
		}
		if userContext.ExpectCookies() {
			basicChain = append(basicChain, newCookiesFiller(userContext.CookiesParams, req))
		}
		if userContext.ExpectHeader() {
			basicChain = append(basicChain, newHeaderFiller(userContext.HeaderParams, req.Header))
		}
		if userContext.ExpectHTTP() {
			basicChain = append(basicChain, newHTTPFiller(userContext.HttpParams, req))
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

const ErrorMessageForMethodNotAllowed = "request resource does not support http method '%s'"

func createNotAllowHandler(method string) router.Handler {
	return newHandler(reflect.ValueOf(func() *response.Response {
		return response.New(fmt.Sprintf(ErrorMessageForMethodNotAllowed, method)).Status(http.StatusMethodNotAllowed)
	}))
}

type optionsHandler struct{}

func (o *optionsHandler) Build(allowMethods string) router.Handler {
	return newHandler(reflect.ValueOf(func() *response.Response {
		return response.New("").
			Headers(response.Headers{
				"Allow": allowMethods,
			})
	}))
}

func (h *handler) getRoute() string {
	return h.route
}

// WildcardIndex implements router.Handler
func (h *handler) WildcardIndex(i int) {
	h.wildcardIndex = i
}
