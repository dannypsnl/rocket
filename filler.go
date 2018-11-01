package rocket

import (
	"encoding/json"
	"io"
	"net/url"
	"reflect"
)

type (
	contextFiller interface {
		fill(context reflect.Value) error
	}
	nextFiller struct {
		next contextFiller
	}
	routeFiller struct {
		nextFiller
		h  *handler
		rs []string
	}
	queryFiller struct {
		nextFiller
		h     *handler
		query url.Values
	}
	jsonFiller struct {
		nextFiller
		h    *handler
		body io.Reader
	}
	formFiller struct {
		nextFiller
		h    *handler
		form url.Values
	}
)

func defaultFillerChain(h *handler, rs []string, body io.Reader, query, form url.Values) contextFiller {
	return newRouteFiller(h, rs,
		newQueryFiller(h, query,
			newJSONFiller(h, body,
				newFormFiller(h, form, nil),
			),
		),
	)
}

func (n *nextFiller) nextFill(ctx reflect.Value) error {
	if n.next != nil {
		return n.next.fill(ctx)
	}
	return nil
}
func (r *routeFiller) fill(ctx reflect.Value) error {
	for idx, route := range r.h.routes {
		if isParameter(route) {
			param := r.rs[len(r.rs)-len(r.h.routes)+idx]
			index := r.h.routeParams[idx]
			v := parseParameter(ctx.Elem().Field(index), param)
			ctx.Elem().Field(index).
				Set(v)
		}
	}
	if r.h.matchedPathIndex != -1 {
		i := r.h.routeParams[r.h.matchedPathIndex]
		ctx.Elem().Field(i).
			Set(reflect.ValueOf(r.h.matchedPath))
	}
	return r.nextFill(ctx)
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
	return q.nextFill(ctx)
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
	return j.nextFill(ctx)
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
	return f.nextFill(ctx)
}

func newRouteFiller(h *handler, rs []string, next contextFiller) contextFiller {
	r := &routeFiller{h: h, rs: rs}
	r.next = next
	return r
}
func newQueryFiller(h *handler, query url.Values, next contextFiller) contextFiller {
	q := &queryFiller{h: h, query: query}
	q.next = next
	return q
}
func newJSONFiller(h *handler, body io.Reader, next contextFiller) contextFiller {
	j := &jsonFiller{h: h, body: body}
	j.next = next
	return j
}
func newFormFiller(h *handler, form url.Values, next contextFiller) contextFiller {
	f := &formFiller{h: h, form: form}
	f.next = next
	return f
}
