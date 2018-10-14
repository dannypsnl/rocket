package rocket

import (
	"bytes"
	"encoding/json"
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

	matchedPath    string
	matchPathIndex int
}

func newHandler() *handler {
	return &handler{
		userDefinedContextOffset: -1,
		cookiesOffset:            -1,
		headerOffset:             -1,
	}
}

func newErrorHandler(code int, content string) *handler {
	h := newHandler()
	h.do = reflect.ValueOf(func() *response.Response {
		return response.New(content).Status(code)
	})
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

func (h *handler) addMatchedPathValueIntoContext(path ...string) {
	buf := bytes.NewBuffer([]byte(``))
	for _, v := range path {
		buf.WriteString(v)
		buf.WriteRune('/')
	}
	h.matchedPath = buf.String()[:buf.Len()-1]
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

func (h *handler) context(rs []string, req *http.Request) []reflect.Value {
	param := make([]reflect.Value, h.do.Type().NumIn())
	if h.hasUserDefinedContext() {
		contextType := h.do.Type().In(h.userDefinedContextOffset)
		context := reflect.New(contextType.Elem())

		if h.expectJsonRequest {
			v := context.Interface()
			err := json.NewDecoder(req.Body).Decode(v)
			if err != nil {
				param[h.userDefinedContextOffset] = reflect.ValueOf(errors.New("400"))
				return param
			}
			param[h.userDefinedContextOffset] = reflect.ValueOf(v)
			return param
		}

		for idx, route := range h.routes {
			if isParameter(route) {
				param := rs[len(rs)-len(h.routes)+idx]
				index := h.routeParams[idx]
				value := parseParameter(context.Elem().Field(index), param)
				context.Elem().Field(index).
					Set(value)
			}
		}

		if h.matchedPath != "" {
			param := h.matchedPath
			index := h.matchPathIndex
			value := parseParameter(context.Elem().Field(index), param)
			context.Elem().Field(index).
				Set(value)
		}

		for k, idx := range h.queryParams {
			values := req.URL.Query()
			field := context.Elem().Field(idx)
			if v, ok := values[k]; ok {
				p := v[0]
				value := parseParameter(context.Elem().Field(idx), p)
				field.Set(value)
			}
		}

		req.ParseForm()
		for k, idx := range h.formParams {
			if v, ok := req.Form[k]; ok {
				p := v[0]
				value := parseParameter(context.Elem().Field(idx), p)
				context.Elem().Field(idx).
					Set(value)
			}
		}

		param[h.userDefinedContextOffset] = context
	}

	if h.needCookies() {
		param[h.cookiesOffset] = reflect.ValueOf(&Cookies{req: req})
	}

	if h.needHeader() {
		param[h.headerOffset] = reflect.ValueOf(&Headers{req: req})
	}

	return param
}
