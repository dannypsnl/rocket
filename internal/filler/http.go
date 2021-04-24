package filler

import (
	"net/http"
	"reflect"
)

type httpFiller struct {
	httpParams map[string]int
	req        *http.Request
}

func NewHTTPFiller(httpParams map[string]int, req *http.Request) Filler {
	return &httpFiller{
		httpParams: httpParams,
		req:        req,
	}
}

func (h *httpFiller) Fill(ctx reflect.Value) error {
	for _, fieldIndex := range h.httpParams {
		field := ctx.Elem().Field(fieldIndex)
		field.Set(reflect.ValueOf(h.req))
	}
	return nil
}
