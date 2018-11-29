package fairing_test

import (
	"net/http"
	"testing"

	asserter "github.com/dannypsnl/assert"

	"github.com/dannypsnl/rocket/fairing"
)

func TestFairingRequestCanModifyRequest(t *testing.T) {
	assert := asserter.NewTester(t)
	hook := fairing.OnRequest(func(r *http.Request) *http.Request {
		r.Header.Set("accept", "application/json")
		r.Cookies()
		r.AddCookie(&http.Cookie{
			Name:  "NEW_HELLO",
			Value: "NEW_WORLD",
		})
		r.Method = "POST"
		return r
	})
	r := &http.Request{
		Method: "GET",
		Header: make(map[string][]string),
	}
	r.Header.Set("accept", "text/html")
	r.AddCookie(&http.Cookie{
		Name:  "HELLO",
		Value: "WORLD",
	})
	r = hook.Invoke(r)

	assert.Eq(r.Method, "POST")
	assert.Eq(r.Header.Get("accept"), "application/json")
	helloCookie, err := r.Cookie("NEW_HELLO")
	assert.NoErr(err)
	assert.Eq(helloCookie.Value, "NEW_WORLD")
}