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
	params            map[int]int // Never custom it. It only for rocket inside.
	formParams        map[string]int
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
				index := h.params[idx]
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
		routes:     strings.Split(*route, "/")[1:],
		do:         handlerDo,
		method:     method,
		params:     make(map[int]int),
		formParams: make(map[string]int),
	}

	handlerT := reflect.TypeOf(do)
	if handlerT.NumIn() > 0 {
		userDefinedT := handlerT.In(0).Elem()
		for idx, r := range h.routes {
			if r[0] == ':' {
				for i := 0; i < userDefinedT.NumField(); i++ {
					key := userDefinedT.Field(i).Tag.Get("route")
					if key == r[1:] {
						h.params[idx] = i
						break
					}
				}
			} else if r[0] == '*' {
				for i := 0; i < userDefinedT.NumField(); i++ {
					key := userDefinedT.Field(i).Tag.Get("route")
					if key == r[1:] {
						h.matchPathIndex = i
						break
					}
				}
			}
		}

		for i := 0; i < userDefinedT.NumField(); i++ {
			key, ok := userDefinedT.Field(i).Tag.Lookup("form")
			if ok {
				h.formParams[key] = i
			}
		}

		for i := 0; i < userDefinedT.NumField(); i++ {
			_, ok := userDefinedT.Field(i).Tag.Lookup("json")
			h.expectJsonRequest = ok
		}
	}

	return h
}

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
