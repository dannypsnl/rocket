package rocket

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/dannypsnl/rocket/internal/context"
	"github.com/dannypsnl/rocket/response"
	"github.com/dannypsnl/rocket/router"
)

type handler struct {
	route  string
	routes []string
	do     reflect.Value // do should return response for HTTP writer
	method string

	guards       []*context.UserContext
	userContexts []*context.UserContext

	matchedPathIndex int
}

func newHandler(do reflect.Value) *handler {
	return &handler{
		do:               do,
		matchedPathIndex: -1,
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

	req.ParseForm()
	for i, userContext := range contexts {
		basicChain := []filler{
			newRouteFiller(
				h.routes,
				reqURL,
				userContext.RouteParams,
				h.matchedPathIndex,
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

const ErrorMessageForMethodNotAllowed = "request resource does not support http method '%s'"

func notAllowHandler(method string) router.Handler {
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

// Routes implements router.Handler
func (h *handler) Route() string {
	return h.route
}

// WildcardIndex implements router.Handler
func (h *handler) WildcardIndex(i int) error {
	h.matchedPathIndex = i
	return nil
}
