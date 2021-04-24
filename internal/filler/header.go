package filler

import (
	"net/http"
	"reflect"

	"github.com/dannypsnl/rocket/internal/parse"
)

type headerFiller struct {
	headerParams map[string]int
	header       http.Header
}

func NewHeaderFiller(headerParams map[string]int, header http.Header) Filler {
	return &headerFiller{
		headerParams: headerParams,
		header:       header,
	}
}

func (h *headerFiller) Fill(ctx reflect.Value) error {
	for key, fieldIndex := range h.headerParams {
		field := ctx.Elem().Field(fieldIndex)
		param := h.header.Get(key)
		value, err := parse.ParseParameter(field.Type(), param)
		if err != nil {
			return err
		}
		field.Set(value)
	}
	return nil
}
