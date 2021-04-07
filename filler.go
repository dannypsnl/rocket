package rocket

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"

	"github.com/dannypsnl/rocket/internal/filepath"
)

type (
	filler interface {
		fill(ctx reflect.Value) error
	}
	routeFiller struct {
		routes                    []string
		routeParams               map[int]int
		reqURL                    []string
		wildcardIndexInSplitRoute int
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
	multiFormFiller struct {
		req             *http.Request
		multiFormParams map[string]int
		limit           map[string]int64
	}
	httpFiller struct {
		httpParams map[string]int
		req        *http.Request
	}
)

func newRouteFiller(routes, reqURL []string, routeParams map[int]int, wildcardIndex int) filler {
	return &routeFiller{
		routes:                    routes,
		routeParams:               routeParams,
		reqURL:                    reqURL,
		wildcardIndexInSplitRoute: wildcardIndex,
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

func newMultiFormFiller(multiFormParams map[string]int, limit map[string]int64, req *http.Request) filler {
	return &multiFormFiller{multiFormParams: multiFormParams, limit: limit, req: req}
}

func newHTTPFiller(httpParams map[string]int, req *http.Request) filler {
	return &httpFiller{
		httpParams: httpParams,
		req:        req,
	}
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
	if r.wildcardIndexInSplitRoute != -1 {
		fieldOffset := r.routeParams[r.wildcardIndexInSplitRoute]
		ctx.Elem().Field(fieldOffset).
			Set(reflect.ValueOf(
				filepath.Join(r.reqURL[r.wildcardIndexInSplitRoute:]...),
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

func (m *multiFormFiller) fill(ctx reflect.Value) error {
	for k, idx := range m.multiFormParams {
		// left shift 20 offset means MB
		err := m.req.ParseMultipartForm(int64(m.limit[k]) << 20)
		if err != nil {
			return err
		}
		file, _, err := m.req.FormFile(k)
		if err != nil {
			return err
		}
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			return err
		}
		field := ctx.Elem().Field(idx)
		value, err := parseParameter(field.Type(), string(fileBytes))
		if err != nil {
			return err
		}
		field.Set(value)
		err = file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *httpFiller) fill(ctx reflect.Value) error {
	for _, fieldIndex := range h.httpParams {
		field := ctx.Elem().Field(fieldIndex)
		field.Set(reflect.ValueOf(h.req))
	}
	return nil
}
