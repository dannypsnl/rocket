package filler

import (
	"encoding/json"
	"io"
	"reflect"
)

type jsonFiller struct {
	body io.Reader
}

func NewJSONFiller(body io.Reader) Filler {
	return &jsonFiller{body: body}
}

func (j *jsonFiller) Fill(ctx reflect.Value) error {
	v := ctx.Interface()
	err := json.NewDecoder(j.body).Decode(v)
	if err != nil {
		return err
	}
	ctx.Elem().Set(reflect.ValueOf(v).Elem())
	return nil
}
