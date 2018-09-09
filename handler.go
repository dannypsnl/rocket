package rocket

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strings"
)

type handler struct {
	routes            []string
	routeParams       map[int]int // Never custom it. It only for rocket inside.
	formParams        map[string]int
	queryParams       map[string]int
	expectJsonRequest bool
	do                reflect.Value // do should return response for HTTP writer
	method            string

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
	if h.do.Type().NumIn() > 0 {
		contextType := h.do.Type().In(0)
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
	return param
}

func handlerByMethod(route *string, do interface{}, method string) *handler {
	handlerDo := reflect.ValueOf(do)
	h := &handler{
		routes:      strings.Split(strings.Trim(*route, "/"), "/"),
		do:          handlerDo,
		method:      method,
		routeParams: make(map[int]int),
		formParams:  make(map[string]int),
		queryParams: make(map[string]int),
	}

	handlerFuncT := reflect.TypeOf(do)
	if handlerFuncT.NumIn() > 0 {
		contextT := handlerFuncT.In(0).Elem()

		routeParams := make(map[string]int)
		for i := 0; i < contextT.NumField(); i++ {
			key, ok := contextT.Field(i).Tag.Lookup("route")
			if ok {
				routeParams[key] = i
			}
		}

		for idx, r := range h.routes {
			// a route part like `:name`
			if r[0] == ':' || r[0] == '*' {
				// r[1:] is `name`, that's the key we expected
				h.routeParams[idx] = routeParams[r[1:]]
			}
		}

		for i := 0; i < contextT.NumField(); i++ {
			key, ok := contextT.Field(i).Tag.Lookup("form")
			if ok {
				h.formParams[key] = i
			}
			key, ok = contextT.Field(i).Tag.Lookup("query")
			if ok {
				h.queryParams[key] = i
			}
			_, ok = contextT.Field(i).Tag.Lookup("json")
			if !h.expectJsonRequest && ok {
				h.expectJsonRequest = ok
			}
		}
	}

	return h
}
