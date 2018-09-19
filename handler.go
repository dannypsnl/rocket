package rocket

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
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

func (h *handler) Handle(rs []string, w http.ResponseWriter, r *http.Request) {
	context := h.context(rs, r)
	response := h.do.Call(
		context,
	)[0].Interface()

	if h.needCookies() {
		for _, c := range context[h.cookiesOffset].Interface().(*Cookies).listOfCookie {
			http.SetCookie(w, c)
		}
	}

	switch response.(type) {
	case *Response:
		res := response.(*Response)
		w.Header().Set("Content-Type", contentTypeOf(res.body))
		for k, v := range res.headers {
			w.Header().Set(k, v)
		}
		fmt.Fprint(w, res.body)
	default:
		w.Header().Set("Content-Type", contentTypeOf(response))
		fmt.Fprint(w, response)
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
			if v, ok := values[k]; ok {
				p := v[0]
				value := parseParameter(context.Elem().Field(idx), p)
				context.Elem().Field(idx).
					Set(value)
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
		param[h.headerOffset] = reflect.ValueOf(&Header{req: req})
	}

	return param
}
