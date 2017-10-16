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

func (rk *Rocket) Dump() {
	fmt.Printf("match: %#v\n", rk.matchs)
	fmt.Printf("handlers: %#v\n", rk.handlers)
}

func Ignite(port string) *Rocket {
	return &Rocket{
		port:     port,
		handlers: make(map[string]Handler),
	}
}

func (rk *Rocket) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var match string
	for _, m := range rk.matchs { // rk.matchs are those static routes
		if m != "/" {
			matched, err := regexp.MatchString(m, r.URL.Path)
			if matched && err == nil {
				match = m
				break
			}
		}
	}
	h := rk.handlers[match]
	matchEls := strings.Split(match, "/")
	Context := make(map[string]string)
	splitRqUrl := strings.Split(r.URL.Path, "/")
	j := 0
	for i, p := range splitRqUrl {
		if matchEls[i] == legalCharsInUrl {
			Context[h.params[j]] = p
			j++
		} else if matchEls[i] == ".*?" {
			Context[h.params[j]] = strings.Join(splitRqUrl[i:], "/")
			break
		}
	}

	fmt.Fprintf(w, h.Do(Context))
}
