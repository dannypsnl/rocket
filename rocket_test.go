package rocket

import (
	"regexp"
	"testing"
)

func SplitContext(t *testing.T, route string, expectedMatch string, lengthOfParamsExpected int, id int) {
	match, params := splitMountUrl(route)
	if match != expectedMatch {
		t.Error(`Number`, id, `match should be '`, expectedMatch, `', but is `, match)
	}
	if len(params) != lengthOfParamsExpected {
		t.Error(`params should have `, len(params), ` params, but it's `, params)
	}
}

func TestSplit(t *testing.T) {
	SplitContext(t, "/", "/", 0, 0)
	SplitContext(t, "/home/:name", "/home/"+legalCharsInUrl, 1, 1)
	SplitContext(t, "/home/:name/age/:age", "/home/"+legalCharsInUrl+"/age/"+legalCharsInUrl, 2, 2)
	SplitContext(t, "/home/:name/:age", "/home/"+legalCharsInUrl+"/"+legalCharsInUrl, 2, 2)
	SplitContext(t, "/home/dan/*name", "/home/dan/.*?", 1, 3)
}

var hello = Get("/:name/age/:age", func(Ctx) Res {
	return "hello"
})

func TestRegex(t *testing.T) {
	rk := Ignite(":8080").
		Mount("/hello", hello)
	r, _ := regexp.Compile(rk.gets[0])
	if !r.MatchString("/hello/dan/age/20") {
		t.Error("Match should success, but it is ", rk.gets[0])
	}
}

func TestVerifyBase(t *testing.T) {
	test_strs := []string{
		"/*path",
		"/hello/:name",
	}
	for _, str := range test_strs {
		func(str string) {
			defer func() {
				if p := recover(); p == nil {
					t.Error("Invalid route didn't panic the program")
				}
			}()
			verifyBase(str)
		}(str)
	}
}

func TestContextType(t *testing.T) {
	ctx := Ctx{"name": "danny"}
	if ctx["name"] != "danny" {
		t.Error("Alias of map should workable, Context>name is ", ctx["name"])
	}
}

func TestFoundMatch(t *testing.T) {
	rk := Ignite(":8080").
		Mount("/", hello)
	_, match, err := rk.foundMatch("/danny/age/20", "GET")
	if err != nil {
		t.Error(`404, but have this handler, bug`)
	}
	//if &h != hello {
	//t.Error(`handler different`)
	//}
	if match != "/"+legalCharsInUrl+"/age/"+legalCharsInUrl {
		t.Error(`We have a incorrect match URL generator`)
	}
}

func TestMethodMatchs(t *testing.T) {
	rk := Ignite(":8080").
		Mount("/", hello)
	if rk.methodMatchs("GET") != &rk.gets {
		t.Error("get fail")
	}
	if rk.methodMatchs("POST") != &rk.posts {
		t.Error("post fail")
	}
	if rk.methodMatchs("PUT") != &rk.puts {
		t.Error("put fail")
	}
	if rk.methodMatchs("DELETE") != &rk.deletes {
		t.Error("delete fail")
	}
}

func TestWrongMatchCausePanic(t *testing.T) {
	rk := Ignite(":8080").
		Mount("/", hello)
	defer func() {
		if p := recover(); p == nil {
			t.Error(`wrong method didn't crash the rocket`)
		}
	}()
	rk.methodMatchs("ADD")
}
