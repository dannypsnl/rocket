package rocket

import "strings"

func splitMountUrl(route string) (string, []string) {
	var match string
	var params []string
	for _, url := range strings.Split(route, "/") {
		if strings.HasPrefix(url, ":") {
			match += "/*"
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
