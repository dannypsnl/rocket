package filler

import (
	"reflect"

	"github.com/dannypsnl/rocket/internal/filepath"
	"github.com/dannypsnl/rocket/internal/parse"
)

type routeFiller struct {
	routes                    []string
	routeParams               map[int]int
	reqURL                    []string
	wildcardIndexInSplitRoute int
}

func NewRouteFiller(routes, reqURL []string, routeParams map[int]int, wildcardIndex int) Filler {
	return &routeFiller{
		routes:                    routes,
		routeParams:               routeParams,
		reqURL:                    reqURL,
		wildcardIndexInSplitRoute: wildcardIndex,
	}
}

func (r *routeFiller) Fill(ctx reflect.Value) error {
	baseRouteLen := len(r.reqURL) - len(r.routes)
	for idx, offset := range r.routeParams {
		param := r.reqURL[baseRouteLen+idx]
		field := ctx.Elem().Field(offset)
		v, err := parse.ParseParameter(field.Type(), param)
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
