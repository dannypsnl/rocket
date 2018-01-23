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

// Rocket is our service.
type Rocket struct {
	port     string
	gets     []string
	posts    []string
	puts     []string
	deletes  []string
	handlers map[string]map[string]handler
}

// Mount add handler into our service.
func (rk *Rocket) Mount(route string, h *handler) *Rocket {
	if !verifyBase(route) {
		panic("Base route can not contain dynamic route.")
	}
	route += h.route
	match, params := splitMountUrl(route)
	h.params = params
	matchs := rk.methodMatchs(h.method)
	*matchs = append(*matchs, match)
	rk.handlers[h.method][match] = *h
	return rk
}

// Launch shoot our service.(start server)
func (rk *Rocket) Launch() {
	sort.Sort(matchArray(rk.gets))
	http.HandleFunc("/", rk.serveLoop)
	log.Fatal(http.ListenAndServe(rk.port, nil))
}

// String serialize info of our service.
func (rk *Rocket) String() string {
	sort.Sort(matchArray(rk.gets))
	return fmt.Sprintf("gets: %#v\n", rk.gets) + fmt.Sprintf("handlers: %#v\n", rk.handlers)
}

// Ignite initial service by port.
func Ignite(port string) *Rocket {
	// TODO: Should use *handler
	hs := make(map[string]map[string]handler)
	// Initial internal method map
	hs["GET"] = make(map[string]handler)
	hs["POST"] = make(map[string]handler)
	hs["PUT"] = make(map[string]handler)
	hs["DELETE"] = make(map[string]handler)
	return &Rocket{
		port:     port,
		handlers: hs,
	}
}

func getContext(h handler, match, path string) map[string]string {
	Context := make(map[string]string)
	matchEls := strings.Split(match, "/")
	splitRqUrl := strings.Split(path, "/")
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
	return Context
}

// serveLoop is prepare for http server trait, but the plan change, it need a new name.
func (rk *Rocket) serveLoop(w http.ResponseWriter, r *http.Request) {
	h, match, err := rk.foundMatch(r.URL.Path, r.Method)
	fmt.Printf("Rquest URL: %#v\n", r.URL.Path)
	if err != nil {
		fmt.Fprintf(w, "404 not found\n")
		return // If 404, we don't need to do others things anymore
	}

	response := h.do(getContext(h, match, r.URL.Path))
	fmt.Fprintf(w, "%s", response)
}

func (rk *Rocket) foundMatch(path string, method string) (handler, string, error) {
	matchs := rk.methodMatchs(method)
	for _, m := range *matchs { // rk.gets are those static routes
		if m == "/" {
			if path == "/" {
				return rk.handlers[method][m], m, nil
			}
		} else {
			matched, err := regexp.MatchString(m, path)
			if matched && err == nil {
				return rk.handlers[method][m], m, nil
			}
		}
	}
	return handler{}, "", errors.New("404")
}

func (rk *Rocket) methodMatchs(method string) *[]string {
	switch method {
	case "GET":
		return &rk.gets
	case "POST":
		return &rk.posts
	case "PUT":
		return &rk.puts
	case "DELETE":
		return &rk.deletes
	default:
		panic("No handle this kind method yet!")
	}
}
