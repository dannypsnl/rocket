package rocket

import (
	"encoding/json"
	"io"
	"net/url"
	"reflect"
)

type (
	filler interface {
		fill(ctx reflect.Value) error
	}
	routeFiller struct {
		routes           []string
		routeParams      map[int]int
		reqURL           []string
		matchedPathIndex int
		matchedPath      string
	}
	queryFiller struct {
		queryParams map[string]int
		query       url.Values
	}
	jsonFiller struct {
		body io.Reader
	}
	formFiller struct {
		formParams map[string]int
		form       url.Values
	}
)

func (r *routeFiller) fill(ctx reflect.Value) error {
	baseRouteLen := len(r.reqURL) - len(r.routes)
	for idx, route := range r.routes {
		if isParameter(route) {
			if index, ok := r.routeParams[idx]; ok {
				param := r.reqURL[baseRouteLen+idx]
				field := ctx.Elem().Field(index)
				v, err := parseParameter(field.Type(), param)
				if err != nil {
					return err
				}
				field.Set(v)
			}
		}
	}
	if r.matchedPathIndex != -1 {
		i := r.routeParams[r.matchedPathIndex]
		ctx.Elem().Field(i).
			Set(reflect.ValueOf(r.matchedPath))
	}
	return nil
}
func (q *queryFiller) fill(ctx reflect.Value) error {
	for k, idx := range q.queryParams {
		field := ctx.Elem().Field(idx)
		if v, ok := q.query[k]; ok {
			param := v[0]
			value, err := parseParameter(field.Type(), param)
			if err != nil {
				return err
			}
			field.Set(value)
		}
	}
	return nil
}
func (j *jsonFiller) fill(ctx reflect.Value) error {
	v := ctx.Interface()
	err := json.NewDecoder(j.body).Decode(v)
	if err != nil {
		return err
	}
	ctx.Elem().Set(reflect.ValueOf(v).Elem())
	return nil
}
func (f *formFiller) fill(ctx reflect.Value) error {
	for k, idx := range f.formParams {
		if v, ok := f.form[k]; ok {
			field := ctx.Elem().Field(idx)
			p := v[0]
			value, err := parseParameter(field.Type(), p)
			if err != nil {
				return err
			}
			field.Set(value)
		}
	}
	return nil
}

type Chain struct {
	ctx reflect.Value
	err error
}

func newChain(ctx reflect.Value) *Chain {
	return &Chain{ctx: ctx}
}

func (c *Chain) pipe(filler filler) *Chain {
	if c.err == nil {
		c.err = filler.fill(c.ctx)
	}
	return c
}

func (c *Chain) error() error {
	return c.err
}

func newRouteFiller(routes, reqURL []string, routeParams map[int]int, matchedPathIndex int, matchedPath string) filler {
	return &routeFiller{
		routes:           routes,
		routeParams:      routeParams,
		reqURL:           reqURL,
		matchedPathIndex: matchedPathIndex,
		matchedPath:      matchedPath,
	}
}
func newQueryFiller(queryParams map[string]int, query url.Values) filler {
	return &queryFiller{
		queryParams: queryParams,
		query:       query,
	}
}
func newJSONFiller(body io.Reader) filler {
	return &jsonFiller{body: body}
}
func newFormFiller(formParams map[string]int, form url.Values) filler {
	return &formFiller{formParams: formParams, form: form}
}
