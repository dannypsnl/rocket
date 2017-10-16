package rocket

import (
	"regexp"
	"testing"
)

func SplitContext(t *testing.T, route string, expectedMatch string, lengthOfParamsExpected int, id int) {
	match, params := splitMountUrl2(route)
	if match != expectedMatch {
		t.Error(`Number`, id, `match should be '`, expectedMatch, `', but is `, match)
	}
	if len(params) != lengthOfParamsExpected {
		t.Error(`params should have `, len(params), ` params, but it's `, params)
	}
}

func TestSplit(t *testing.T) {
	SplitContext(t, "/", "/", 0, 0)
	SplitContext(t, "/home/:name", "/home/*", 1, 1)
	SplitContext(t, "/home/:name/age/:age", "/home/*/age/*", 2, 2)
	SplitContext(t, "/home/:name/:age", "/home/*/*", 2, 2)
	SplitContext(t, "/home/dan/*name", "/home/dan/.*?", 1, 3)
}

func TestRegex(t *testing.T) {
	legalCharsInUrl := "[a-zA-Z0-9-_]+"
	r, _ := regexp.Compile("/home/" + legalCharsInUrl + "/src")
	r2, _ := regexp.Compile("/home/*/src")
	if !r.MatchString("/home/dan/src") && !r2.MatchString("/home/dan/src") {
		t.Error("fail")
	}
	if r.MatchString("/home/dan/20/src") && r2.MatchString("/home/dan/20/src") {
		t.Error("fail")
	}
}
