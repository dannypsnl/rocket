package rocket

import (
	"encoding/json"
	"io"
	"net/http"
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
	}
	queryFiller struct {
		queryParams map[string]int
		query       url.Values
	}
	headerFiller struct {
		headerParams map[string]int
		header       http.Header
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

func newRouteFiller(routes, reqURL []string, routeParams map[int]int, matchedPathIndex int) filler {
	return &routeFiller{
		routes:           routes,
		routeParams:      routeParams,
		reqURL:           reqURL,
		matchedPathIndex: matchedPathIndex,
	}
}

func newQueryFiller(queryParams map[string]int, query url.Values) filler {
	return &queryFiller{
		queryParams: queryParams,
		query:       query,
	}
}

func newHeaderFiller(headerParams map[string]int, header http.Header) filler {
	return &headerFiller{
		headerParams: headerParams,
		header:       header,
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
		paths := r.reqURL[r.matchedPathIndex:]

		lenOfPath := len(paths[0])
		for _, p := range paths[1:] {
			lenOfPath += len(p) + 1
		}
		path := make([]byte, lenOfPath)
		lastOne := len(paths) - 1
		index := 0
		for _, v := range paths[:lastOne] {
			copy(path[index:], []byte(v))
			index += len(v)
			copy(path[index:], []byte{'/'})
			index++
		}
		v := paths[lastOne]
		copy(path[index:], []byte(v))

		ctx.Elem().Field(i).
			Set(reflect.ValueOf(
				string(path),
			))
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

func (h *headerFiller) fill(ctx reflect.Value) error {
	for key, fieldIndex := range h.headerParams {
		field := ctx.Elem().Field(fieldIndex)
		param := h.header.Get(key)
		value, err := parseParameter(field.Type(), param)
		if err != nil {
			return err
		}
		field.Set(value)
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
