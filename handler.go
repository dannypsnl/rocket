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
		basicChain := newPipeline().
			pipe(newRouteFiller(
				h.routes,
				reqURL,
				userContext.RouteParams,
				h.matchedPathIndex,
				h.matchedPath,
			)).
			pipe(newQueryFiller(userContext.QueryParams, req.URL.Query()))
		if userContext.IsHeaders {
			userContexts[i] = reflect.ValueOf(&Headers{header: req.Header})
			return userContexts, nil
		} else if userContext.ExpectJSONRequest {
			basicChain.pipe(newJSONFiller(req.Body))
		} else {
			basicChain.pipe(newFormFiller(userContext.FormParams, req.Form))
		}
		if userContext.ExpectCookies() {
			basicChain.pipe(newCookiesFiller(userContext.CookiesParams, req))
		}
		ctx, err := basicChain.
			run(userContext)
		if err != nil {
			return nil, err
		}
		userContexts[i] = ctx
	}
	return userContexts, nil
}
