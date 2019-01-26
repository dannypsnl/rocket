package rocket

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"reflect"

	"github.com/dannypsnl/rocket/internal/context"
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
	cookiesFiller struct {
		cookiesParams map[string]int
		req           *http.Request
	}
	jsonFiller struct {
		body io.Reader
	}
	formFiller struct {
		formParams map[string]int
		form       url.Values
	}
)

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

func newCookiesFiller(cookiesParams map[string]int, req *http.Request) filler {
	return &cookiesFiller{
		cookiesParams: cookiesParams,
		req:           req,
	}
}

func newJSONFiller(body io.Reader) filler {
	return &jsonFiller{body: body}
}

func newFormFiller(formParams map[string]int, form url.Values) filler {
	return &formFiller{formParams: formParams, form: form}
}

func (r *routeFiller) fill(ctx reflect.Value) error {
	baseRouteLen := len(r.reqURL) - len(r.routes)
	for idx, offset := range r.routeParams {
		param := r.reqURL[baseRouteLen+idx]
		field := ctx.Elem().Field(offset)
		v, err := parseParameter(field.Type(), param)
		if err != nil {
			return err
		}
		field.Set(v)
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

func (c *cookiesFiller) fill(ctx reflect.Value) error {
	for key, fieldIndex := range c.cookiesParams {
		field := ctx.Elem().Field(fieldIndex)
		cookie, err := c.req.Cookie(key)
		if err != nil {
			return err
		}
		field.Set(reflect.ValueOf(cookie))
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

type pipeline struct {
	fillers []filler
}

func newPipeline() *pipeline {
	return &pipeline{
		fillers: make([]filler, 0),
	}
}

func (p *pipeline) pipe(f filler) *pipeline {
	p.fillers = append(p.fillers, f)
	return p
}

func (p *pipeline) run(userContext *context.UserContext) (reflect.Value, error) {
	ctx := reflect.New(userContext.ContextType)
	for _, filler := range p.fillers {
		if err := filler.fill(ctx); err != nil {
			return ctx, err
		}
	}
	return ctx, nil
}
