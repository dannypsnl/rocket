package rocket

func splitMountUrl(route string) (string, []string) {
	match := ""

	firstTime := true
	open := false
	neveropen := true
	start := 0

	var params []string
	for i, r := range route {
		if r == ':' {
			if firstTime {
				match = route[:i-1]
				firstTime = false
			}
			start = i + 1
			open = true
			neveropen = false
		}
		if r == '*' {
			match = route[:i-1] + "/.*?"
			params = append(params, route[i:])
			break
		}
		if i == len(route)-1 {
			if neveropen {
				match = route
			} else {
				match += "/*"
				params = append(params, route[start:i+1])
			}
		}
		if open && r == '/' {
			// Get param setting string.
			match += "/*"
			if i != len(route)-1 {
				match += "/"
			}
			params = append(params, route[start:i])
			open = false
		}
		if !open && r != '/' {
			match += string(route[i])
		}
	}
	return match, params
}
