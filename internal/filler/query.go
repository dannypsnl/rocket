package filler

import (
	"net/url"
	"reflect"

	"github.com/dannypsnl/rocket/internal/parse"
)

type queryFiller struct {
	queryParams map[string]int
	query       url.Values
}

func NewQueryFiller(queryParams map[string]int, query url.Values) Filler {
	return &queryFiller{
		queryParams: queryParams,
		query:       query,
	}
}

func (q *queryFiller) Fill(ctx reflect.Value) error {
	for k, idx := range q.queryParams {
		field := ctx.Elem().Field(idx)
		if v, ok := q.query[k]; ok {
			param := v[0]
			value, err := parse.ParseParameter(field.Type(), param)
			if err != nil {
				return err
			}

			field.Set(value)
		}
	}
	return nil
}
