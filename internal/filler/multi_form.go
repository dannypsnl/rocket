package filler

import (
	"net/http"
	"reflect"

	"github.com/dannypsnl/rocket/internal/parse"
)

type multiFormFiller struct {
	req             *http.Request
	multiFormParams map[string]int
	paramIsFile     map[string]bool
	limit           int64
}

func NewMultiFormFiller(limit int64, multiFormParams map[string]int, paramIsFile map[string]bool, req *http.Request) Filler {
	return &multiFormFiller{limit: limit, multiFormParams: multiFormParams, paramIsFile: paramIsFile, req: req}
}

func (m *multiFormFiller) Fill(ctx reflect.Value) error {
	err := m.req.ParseMultipartForm(m.limit << 20)
	if err != nil {
		return err
	}
	for k, idx := range m.multiFormParams {
		field := ctx.Elem().Field(idx)
		if m.paramIsFile[k] {
			file, _, err := m.req.FormFile(k)
			if err != nil {
				return err
			}
			field.Set(reflect.ValueOf(file))
		} else {
			v := m.req.FormValue(k)
			value, err := parse.ParseParameter(field.Type(), string(v))
			if err != nil {
				return err
			}
			field.Set(value)
		}
	}
	return nil
}
