package rocket

import (
	"regexp"
)

const legalCharsInUrl = "([a-zA-Z0-9-_]+)"

func verifyBase(route string) bool {
	r, _ := regexp.Compile(".*?[:*].*?")
	// Contains : part will Match, it can be on a Base Route
	if r.MatchString(route) {
		panic("Base route can not contain dynamic route.")
	}
	return true
}
