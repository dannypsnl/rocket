package rocket

import (
	"errors"
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
	gets     []string
	handlers map[string]handler
}

func (r *Rocket) Mount(route string, h *handler) *Rocket {
	if !verifyBase(route) {
		panic("Base route can not contain dynamic route.")
	}
	route += h.route
	match, params := splitMountUrl(route)
	h.params = params
	r.gets = append(r.gets, match)
	r.handlers[match] = *h
	return r
}

func (rk *Rocket) Launch() {
	sort.Sort(matchArray(rk.gets))
	http.HandleFunc("/", rk.ServeHTTP)
	log.Fatal(http.ListenAndServe(rk.port, nil))
}

func (rk *Rocket) Dump() {
	sort.Sort(matchArray(rk.gets))
	fmt.Printf("gets: %#v\n", rk.gets)
	fmt.Printf("handlers: %#v\n", rk.handlers)
}

func Ignite(port string) *Rocket {
	return &Rocket{
		port:     port,
		handlers: make(map[string]handler),
	}
}
func (rk *Rocket) foundMatch(r *http.Request) (string, error) {
	for _, m := range rk.gets { // rk.gets are those static routes
		if m == "/" {
			if r.URL.Path == "/" {
				return m, nil
			}
		} else {
			matched, err := regexp.MatchString(m, r.URL.Path)
			if matched && err == nil {
				return m, nil
			}
		}
	}
	return "", errors.New("404")
}
func (rk *Rocket) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	match, err := rk.foundMatch(r)
	if err != nil {
		fmt.Fprintf(w, "404 not found\n")
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
	response := h.do(Context)
	fmt.Fprintf(w, "%s", response)
}
