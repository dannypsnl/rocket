package rocket

import (
	"encoding/json"
	"io"
	"net/url"
	"reflect"
)

type (
	filler interface {
		fill(context reflect.Value) error
	}
	getNextFiller struct {
		nextFiller filler
	}
	routeFiller struct {
		getNextFiller
		h      *handler
		reqURL []string
	}
	queryFiller struct {
		getNextFiller
		h     *handler
		query url.Values
	}
	jsonFiller struct {
		getNextFiller
		h    *handler
		body io.Reader
	}
	formFiller struct {
		getNextFiller
		h    *handler
		form url.Values
	}
)

func defaultFillerChain(h *handler, rs []string, body io.Reader, query, form url.Values) filler {
	return newRouteFiller(h, rs,
		newQueryFiller(h, query,
			newJSONFiller(h, body,
				newFormFiller(h, form, nil),
			),
		),
	)
}

func (n *getNextFiller) next(ctx reflect.Value) error {
	if n.nextFiller != nil {
		return n.nextFiller.fill(ctx)
	}
	return nil
}
func (r *routeFiller) fill(ctx reflect.Value) error {
	baseRouteLen := len(r.reqURL) - len(r.h.routes)
	for idx, route := range r.h.routes {
		if isParameter(route) {
			param := r.reqURL[baseRouteLen+idx]
			index := r.h.routeParams[idx]
			field := ctx.Elem().Field(index)
			v := parseParameter(field, param)
			field.Set(v)
		}
	}
	if r.h.matchedPathIndex != -1 {
		i := r.h.routeParams[r.h.matchedPathIndex]
		ctx.Elem().Field(i).
			Set(reflect.ValueOf(r.h.matchedPath))
	}
	return r.next(ctx)
}
func (q *queryFiller) fill(ctx reflect.Value) error {
	for k, idx := range q.h.queryParams {
		field := ctx.Elem().Field(idx)
		if v, ok := q.query[k]; ok {
			p := v[0]
			value := parseParameter(field, p)
			field.Set(value)
		}
	}
	return q.next(ctx)
}
func (j *jsonFiller) fill(ctx reflect.Value) error {
	if j.h.expectJsonRequest {
		v := ctx.Interface()
		err := json.NewDecoder(j.body).Decode(v)
		if err != nil {
			return err
		}
		ctx.Elem().Set(reflect.ValueOf(v).Elem())
		return nil
	}
	return j.next(ctx)
}
func (f *formFiller) fill(ctx reflect.Value) error {
	for k, idx := range f.h.formParams {
		if v, ok := f.form[k]; ok {
			field := ctx.Elem().Field(idx)
			p := v[0]
			value := parseParameter(field, p)
			field.Set(value)
		}
	}
	return f.next(ctx)
}

func newRouteFiller(h *handler, reqURL []string, next filler) filler {
	r := &routeFiller{h: h, reqURL: reqURL}
	r.nextFiller = next
	return r
}
func newQueryFiller(h *handler, query url.Values, next filler) filler {
	q := &queryFiller{h: h, query: query}
	q.nextFiller = next
	return q
}
func newJSONFiller(h *handler, body io.Reader, next filler) filler {
	j := &jsonFiller{h: h, body: body}
	j.nextFiller = next
	return j
}
func newFormFiller(h *handler, form url.Values, next filler) filler {
	f := &formFiller{h: h, form: form}
	f.nextFiller = next
	return f
}
