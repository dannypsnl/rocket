package rocket

import (
	"regexp"
	"strings"
)

const legalCharsInUrl = "([a-zA-Z0-9-_]+)"

func splitMountUrl(route string) (string, []string) {
	var match string
	var params []string
	// '/:id' is params in url.
	// '/*filepath' is params about filepath.
	// '/home, data' is params from post method.
	for _, url := range strings.Split(route, "/") {
		if strings.HasPrefix(url, ":") {
			match += "/" + legalCharsInUrl
			params = append(params, url[1:])
		} else if strings.HasPrefix(url, "*") {
			match += "/.*?"
			params = append(params, url[1:])
			break
		} else if url != "" {
			match += "/" + url
		}
	}
	if match == "" {
		match = "/"
	}
	return match, params
}

func verifyBase(route string) bool {
	r, _ := regexp.Compile(".*?[:*].*?")
	// Contains : part will Match, it can be on a Base Route
	if r.MatchString(route) {
		panic("Base route can not contain dynamic route.")
	}
	return true
}
