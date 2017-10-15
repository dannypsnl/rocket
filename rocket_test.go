package rocket

import "testing"

func SplitContext(t *testing.T, route string, expectedMatch string, lengthOfParamsExpected int) {
	match, params := split(route)
	if match != expectedMatch {
		t.Error(`match should be '/', but is `, match)
	}
	if len(params) != lengthOfParamsExpected {
		t.Error(`params should have `, len(params), ` params, but it's `, params)
	}
}

func TestSplit(t *testing.T) {
	SplitContext(t, "/", "/", 0)
	SplitContext(t, "/home/:name", "/home", 1)
	SplitContext(t, "/home/:name/age/:age", "/home", 2)
	SplitContext(t, "/home/dan/*name", "/home/dan", 1)
}
