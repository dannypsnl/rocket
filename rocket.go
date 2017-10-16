package rocket

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type Rocket struct {
	port     string
	matchs   []string
	handlers map[string]Handler
}

func (r *Rocket) Mount(route string, h Handler) *Rocket {
	route += h.Route
	// '/:id' is params in url.
	// '/*filepath' is params about filepath.
	// '/home, data' is params from post method.
	match, params := splitMountUrl(route)
	h.params = params
	r.matchs = append(r.matchs, match)
	r.handlers[match] = h
	return r
}

func (rk *Rocket) Launch() {
	http.HandleFunc("/", rk.ServeHTTP)
	log.Fatal(http.ListenAndServe(rk.port, nil))
}

func Ignite(port string) *Rocket {
	return &Rocket{
		port:     port,
		handlers: make(map[string]Handler),
	}
}

func (rk *Rocket) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var match string
	var rqUrl string
	for _, m := range rk.matchs { // rk.matchs are those static routes
		reg, _ := regexp.Compile(m)
		if reg.MatchString(r.URL.Path) {
			match = m
			rqUrl = r.URL.Path
			break
		}
	}
	h := rk.handlers[match]
	matchEls := strings.Split(match, "/")
	Context := make(map[string]string)
	splitRqUrl := strings.Split(rqUrl, "/")
	fmt.Fprintf(w, "%#v\t%#v\t%#v%#v\n", matchEls, splitRqUrl, rqUrl, match)
	j := 0
	for i, p := range splitRqUrl {
		if matchEls[i] == "*" {
			Context[h.params[j]] = p
			j++
		} else if matchEls[i] == ".*?" {
			Context[h.params[j]] = strings.Join(splitRqUrl[i:], "/")
			break
		}
	}
	fmt.Fprintf(w, "%#v\t%#v\n", Context, h.params)

	fmt.Fprintf(w, h.Do(Context))
}
