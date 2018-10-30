package rocket

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
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
		do:                       do,
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

func (h *handler) Handle(rs []string, r *http.Request) *response.Response {
	resp := h.do.Call(
		h.context(rs, r),
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

func (h *handler) fillRoute(context reflect.Value, rs []string) {
	for idx, route := range h.routes {
		if isParameter(route) {
			param := rs[len(rs)-len(h.routes)+idx]
			index := h.routeParams[idx]
			value := parseParameter(context.Elem().Field(index), param)
			context.Elem().Field(index).
				Set(value)
		}
	}
	if h.matchedPathIndex != -1 {
		index := h.routeParams[h.matchedPathIndex]
		context.Elem().Field(index).
			Set(reflect.ValueOf(h.matchedPath))
	}
}
func (h *handler) fillQuery(context reflect.Value, querys url.Values) {
	for k, idx := range h.queryParams {
		field := context.Elem().Field(idx)
		if v, ok := querys[k]; ok {
			p := v[0]
			value := parseParameter(context.Elem().Field(idx), p)
			field.Set(value)
		}
	}
}
func (h *handler) fillJSON(context reflect.Value, req *http.Request) (reflect.Value, error) {
	if h.expectJsonRequest {
		v := context.Interface()
		err := json.NewDecoder(req.Body).Decode(v)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(v), nil
	}
	return context, nil
}
func (h *handler) fillForm(context reflect.Value, form url.Values) {
	for k, idx := range h.formParams {
		if v, ok := form[k]; ok {
			p := v[0]
			value := parseParameter(context.Elem().Field(idx), p)
			context.Elem().Field(idx).
				Set(value)
		}
	}
}

type contextParser interface {
	parse(context reflect.Value, rs []string, vs url.Values)
}

func (h *handler) context(rs []string, req *http.Request) []reflect.Value {
	param := make([]reflect.Value, h.do.Type().NumIn())
	if h.hasUserDefinedContext() {
		contextType := h.do.Type().In(h.userDefinedContextOffset).Elem()
		context := reflect.New(contextType)

		h.fillRoute(context, rs)
		h.fillQuery(context, req.URL.Query())
		v, err := h.fillJSON(context, req)
		if err != nil {
			param[h.userDefinedContextOffset] = reflect.ValueOf(errors.New("400"))
			return param
		}
		param[h.userDefinedContextOffset] = v
		req.ParseForm()
		h.fillForm(context, req.Form)

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
