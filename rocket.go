package rocket

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strings"
)

type matchArray []string

func (ma matchArray) Len() int      { return len(ma) }
func (ma matchArray) Swap(i, j int) { ma[i], ma[j] = ma[j], ma[i] }
func (ma matchArray) Less(i, j int) bool {
	return len(strings.Split(ma[i], "/")) > len(strings.Split(ma[j], "/"))
}

type Rocket struct {
	port     string
	matchs   []string
	handlers map[string]Handler
}

func (r *Rocket) Mount(route string, h Handler) *Rocket {
	if !verifyBase(route) {
		panic("Base route can not contain dynamic route.")
	}
	route += h.Route
	match, params := splitMountUrl(route)
	h.params = params
	r.matchs = append(r.matchs, match)
	r.handlers[match] = h
	return r
}

func (rk *Rocket) Launch() {
	sort.Sort(matchArray(rk.matchs))
	http.HandleFunc("/", rk.ServeHTTP)
	log.Fatal(http.ListenAndServe(rk.port, nil))
}

func (rk *Rocket) Dump() {
	sort.Sort(matchArray(rk.matchs))
	fmt.Printf("matchs: %#v\n", rk.matchs)
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
		if m == "/" {
			if r.URL.Path == "/" {
				match = m
				break
			}
		} else {
			matched, err := regexp.MatchString(m, r.URL.Path)
			if matched && err == nil {
				match = m
				break
			}
		}
	}
	fmt.Printf("Rquest URL: %#v\n", r.URL.Path)
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

	// So, at next API, we will get rocket.Response object.
	// TODO: resolve rocket.Response type.
	fmt.Fprintf(w, h.Do(Context))
}
