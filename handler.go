package rocket

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/dannypsnl/rocket/internal/context"
	"github.com/dannypsnl/rocket/internal/filler"
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

	userContexts []*context.UserContext

	wildcardIndex int
	// get config
	rocket *Rocket
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
	contexts, err := h.getFilledContexts(h.userContexts, reqURL, r)
	if err != nil {
		return response.New(err.Error()).
			Status(http.StatusBadRequest)
	}

	if resp := h.verify(contexts); resp != nil {
		return resp
	}

	resp := h.do.Call(
		contexts,
	)[0].Interface()

	switch v := resp.(type) {
	case *response.Response:
		return v
	default:
		return response.New(v)
	}
}

func (h *handler) verify(contexts []reflect.Value) *response.Response {
	for _, c := range contexts {
		switch guard := c.Interface().(type) {
		case Guard:
			if resp := guard.VerifyRequest(); resp != nil {
				return resp
			}
		}
	}
	return nil
}

func (h *handler) getFilledContexts(contexts []*context.UserContext, reqURL []string, req *http.Request) ([]reflect.Value, error) {
	userContexts := make([]reflect.Value, len(contexts))

	err := req.ParseForm()
	if err != nil {
		return nil, err
	}
	for i, userContext := range contexts {
		basicChain := []filler.Filler{
			filler.NewRouteFiller(
				h.routes,
				reqURL,
				userContext.RouteParams,
				h.wildcardIndex,
			),
			filler.NewQueryFiller(userContext.QueryParams, req.URL.Query()),
		}
		if userContext.ExpectJSONRequest {
			basicChain = append(basicChain, filler.NewJSONFiller(req.Body))
		} else if userContext.ExpectMultiFormsRequest {
			basicChain = append(basicChain, filler.NewMultiFormFiller(h.rocket.MultiFormBodySizeLimit, userContext.MultiFormParams, userContext.MultiFormParamsIsFile, req))
		} else {
			basicChain = append(basicChain, filler.NewFormFiller(userContext.FormParams, req.Form))
		}
		if userContext.ExpectCookies() {
			basicChain = append(basicChain, filler.NewCookiesFiller(userContext.CookiesParams, req))
		}
		if userContext.ExpectHeader() {
			basicChain = append(basicChain, filler.NewHeaderFiller(userContext.HeaderParams, req.Header))
		}
		if userContext.ExpectHTTP() {
			basicChain = append(basicChain, filler.NewHTTPFiller(userContext.HttpParams, req))
		}
		ctx := reflect.New(userContext.ContextType)
		for _, f := range basicChain {
			if err := f.Fill(ctx); err != nil {
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
