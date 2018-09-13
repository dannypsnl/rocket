package rocket

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
)

type handler struct {
	routes []string
	do     reflect.Value // do should return response for HTTP writer
	method string

	userDefinedContextOffset int
	cookiesOffset            int

	routeParams       map[int]int // Never custom it. It only for rocket inside.
	formParams        map[string]int
	queryParams       map[string]int
	expectJsonRequest bool

	matchedPath    string
	matchPathIndex int
}

func (h *handler) addMatchedPathValueIntoContext(path ...string) {
	buf := bytes.NewBuffer([]byte(``))
	for _, v := range path {
		buf.WriteString(v)
		buf.WriteRune('/')
	}
	h.matchedPath = buf.String()[:buf.Len()-1]
}

func (h *handler) context(rs []string, req *http.Request) []reflect.Value {
	param := make([]reflect.Value, 0)
	if h.userDefinedContextOffset != -1 {
		contextType := h.do.Type().In(h.userDefinedContextOffset)
		context := reflect.New(contextType.Elem())

		if h.expectJsonRequest {
			v := context.Interface()
			err := json.NewDecoder(req.Body).Decode(v)
			if err != nil {
				param = append(param, reflect.ValueOf(errors.New("400")))
				return param
			}
			param = append(param, reflect.ValueOf(v))
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

		param = append(param, context)
	}

	if h.cookiesOffset != -1 {
		cs := &Cookies{
			req: req,
		}
		param = append(param, reflect.ValueOf(cs))
	}

	return param
}
