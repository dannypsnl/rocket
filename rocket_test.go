package rocket

import (
	"regexp"
	"testing"
)

func SplitContext(t *testing.T, route string, expectedMatch string, lengthOfParamsExpected int, id int) {
	match, params := splitMountUrl(route)
	if match != expectedMatch {
		t.Error(`Number`, id, `match should be '`, expectedMatch, `', but is `, match)
	}
	if len(params) != lengthOfParamsExpected {
		t.Error(`params should have `, len(params), ` params, but it's `, params)
	}
}

func TestSplit(t *testing.T) {
	SplitContext(t, "/", "/", 0, 0)
	SplitContext(t, "/home/:name", "/home/"+legalCharsInUrl, 1, 1)
	SplitContext(t, "/home/:name/age/:age", "/home/"+legalCharsInUrl+"/age/"+legalCharsInUrl, 2, 2)
	SplitContext(t, "/home/:name/:age", "/home/"+legalCharsInUrl+"/"+legalCharsInUrl, 2, 2)
	SplitContext(t, "/home/dan/*name", "/home/dan/.*?", 1, 3)
}

var hello = Handler{
	Route: "/:name/age/:age",
	Do: func(Context) string {
		return "hello"
	},
}

func TestRegex(t *testing.T) {
	rk := Ignite(":8080").
		Mount("/hello", hello)
	r, _ := regexp.Compile(rk.matchs[0])
	if !r.MatchString("/hello/dan/age/20") {
		t.Error("Match should success, but it is ", rk.matchs[0])
	}
}

func TestVerifyBase(t *testing.T) {
	test_strs := []string{
		"/*path",
		"/hello/:name",
	}
	for _, str := range test_strs {
		if verifyBase(str) {
			t.Error("Base route should not contain dynamic part.")
		}
	}
}

func TestContextType(t *testing.T) {
	ctx := Context{"name": "danny"}
	if ctx["name"] != "danny" {
		t.Error("Alias of map should workable, Context>name is ", ctx["name"])
	}
}

type response struct {
	contentType string
	messages    []string
}

type ResponseBuilder struct {
	contentType string
	messages    []string
}

func (rb *ResponseBuilder) Done() response {
	defer func() {
		rb.contentType = ""
		rb.messages = []string{}
	}()
	return response{
		contentType: rb.contentType,
		messages:    rb.messages,
	}
}

func (rb *ResponseBuilder) ContentType(contentType string) *ResponseBuilder {
	rb.contentType = contentType
	return rb
}

func TestResponseBuilder(t *testing.T) {
	rb := ResponseBuilder{}
	if rb.Done().contentType != "" {
		t.Error("ResponseBuilder should had nothing.")
	}
	rb.ContentType("brbrbr")

	if rb.Done().contentType != "brbrbr" {
		t.Error("ContentType should be brbrbr, but is ", rb.Done().contentType)
	}
}
