package filler

import (
	"net/http"
	"reflect"
)

type cookiesFiller struct {
	cookiesParams map[string]int
	req           *http.Request
}

func NewCookiesFiller(cookiesParams map[string]int, req *http.Request) Filler {
	return &cookiesFiller{
		cookiesParams: cookiesParams,
		req:           req,
	}
}

func (c *cookiesFiller) Fill(ctx reflect.Value) error {
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
