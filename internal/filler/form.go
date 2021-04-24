package filler

import (
	"net/url"
	"reflect"

	"github.com/dannypsnl/rocket/internal/parse"
)

type formFiller struct {
	formParams map[string]int
	form       url.Values
}

func NewFormFiller(formParams map[string]int, form url.Values) Filler {
	return &formFiller{formParams: formParams, form: form}
}

func (f *formFiller) Fill(ctx reflect.Value) error {
	for k, idx := range f.formParams {
		if v, ok := f.form[k]; ok {
			field := ctx.Elem().Field(idx)
			p := v[0]
			value, err := parse.ParseParameter(field.Type(), p)
			if err != nil {
				return err
			}
			field.Set(value)
		}
	}
	return nil
}
