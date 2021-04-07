package context

import (
	"fmt"
	"net/http"
	"reflect"
)

type UserContext struct {
	ContextType           reflect.Type
	IsHeaders             bool
	RouteParams           map[int]int
	FormParams            map[string]int
	MultiFormParams       map[string]int
	MultiFormParamsIsFile map[string]bool
	QueryParams           map[string]int
	// `cookie:"token"`, would store "token" as key, field index as value
	CookiesParams map[string]int
	// `header:"Content-Type"`, would store "Content-Type" as key, field index as value
	HeaderParams map[string]int
	// `http:"request"` would take request
	// http tag is limited, can only take what we allow at here:
	//
	// - request
	HttpParams        map[string]int
	ExpectJSONRequest bool
}

func NewUserContext() *UserContext {
	return &UserContext{
		IsHeaders:             false,
		RouteParams:           make(map[int]int),
		FormParams:            make(map[string]int),
		MultiFormParams:       make(map[string]int),
		MultiFormParamsIsFile: make(map[string]bool),
		QueryParams:           make(map[string]int),
		CookiesParams:         make(map[string]int),
		HeaderParams:          make(map[string]int),
		HttpParams:            make(map[string]int),
		ExpectJSONRequest:     false,
	}
}

func (ctx *UserContext) CacheParamsOffset(contextT reflect.Type, routes []string) *UserContext {
	ctx.ContextType = contextT
	routeParams := make(map[string]int)
	for i := 0; i < contextT.NumField(); i++ {
		tagOfField := contextT.Field(i).Tag
		key, ok := tagOfField.Lookup("route")
		if ok {
			routeParams[key] = i
		}
		key, ok = tagOfField.Lookup("form")
		if ok {
			ctx.FormParams[key] = i
		}
		key, ok = tagOfField.Lookup("multiform")
		if ok {
			ctx.MultiFormParams[key] = i
			// means this is a file
			_, ok = tagOfField.Lookup("file")
			if ok {
				ctx.MultiFormParamsIsFile[key] = true
			}
		}
		key, ok = tagOfField.Lookup("query")
		if ok {
			ctx.QueryParams[key] = i
		}
		key, ok = tagOfField.Lookup("cookie")
		if ok {
			if !contextT.Field(i).Type.AssignableTo(reflect.TypeOf(&http.Cookie{})) {
				panic("type of fields those try to be a cookie must be `*http.Cookie`")
			}
			ctx.CookiesParams[key] = i
		}
		key, ok = tagOfField.Lookup("header")
		if ok {
			ctx.HeaderParams[key] = i
		}
		key, ok = tagOfField.Lookup("http")
		if ok {
			switch key {
			case "request":
			default:
				panic(fmt.Sprintf("unknown resource be required in http tag: `%s`", key))
			}
			ctx.HttpParams[key] = i
		}
		// we found json tag, then must expect a JSON request
		_, ok = tagOfField.Lookup("json")
		if !ctx.ExpectJSONRequest && ok {
			ctx.ExpectJSONRequest = ok
		}
	}

	for idx, r := range routes {
		if r[0] == ':' || r[0] == '*' {
			// what if r is ":name", r[1:] is "name", that's the key we expected
			param := r[1:]
			if _, ok := routeParams[param]; ok {
				ctx.RouteParams[idx] = routeParams[param]
			}
		}
	}

	return ctx
}

func (ctx *UserContext) ExpectCookies() bool {
	return len(ctx.CookiesParams) > 0
}

func (ctx *UserContext) ExpectHeader() bool {
	return len(ctx.HeaderParams) > 0
}

func (ctx *UserContext) ExpectHTTP() bool {
	return len(ctx.HttpParams) > 0
}
