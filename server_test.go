package rocket_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/dannypsnl/assert"
	"github.com/dannypsnl/rocket"
)

type User struct {
	Name string `route:"name"`
	Age  int    `route:"age"`

	QueryName string `query:"name"`
}

type ForPost struct {
	Val string `json:"value"`
}

type ForPatch struct {
	Val string `form:"value"`
}

type Files struct {
	FileName string `route:"filename"`
}

var (
	homePage = rocket.Get("/", func() rocket.Html {
		return `
		<h1>Title</h1>
		<p>Hello, World</p>
		`
	})
	staticFiles = rocket.Get("/static/*filename", func(fs *Files) string {
		return fs.FileName
	})
	noParamNoRoute = rocket.Get("/", func() string {
		return "no param no route"
	})
	helloName = rocket.Get("/:name", func(u *User) string {
		return "Hello, " + u.Name
	})
	forPost = rocket.Post("/post", func(f *ForPost) rocket.Json {
		return `{"value": "response"}`
	})
	forPatch = rocket.Patch("/patch", func(f *ForPatch) string {
		return f.Val
	})
	user = rocket.Get("/:name/name", func(u *User) string {
		return u.Name
	})
	query = rocket.Get("/query", func(u *User) string {
		return u.QueryName
	})
)

func TestServer(t *testing.T) {
	assert := assert.NewTester(t)

	rk := rocket.Ignite(":8080").
		Mount("/", homePage, staticFiles).
		Mount("/hello", helloName).
		Mount("/users", user).
		Mount("/test", query, noParamNoRoute, forPatch, forPost).
		Default(func() rocket.Html {
			return "<h1>Page Not Found</h1>"
		})
	ts := httptest.NewServer(rk)
	defer ts.Close()

	t.Run("GetUserName", func(t *testing.T) {
		result, _, err := response("GET", ts.URL, "/users/Danny/name")
		assert.Eq(err, nil)
		assert.Eq(result, "Danny")
	})

	t.Run("GetHomePage", func(t *testing.T) {
		result, header, err := response("GET", ts.URL, "/")
		assert.Eq(err, nil)
		assert.Eq(result, `
		<h1>Title</h1>
		<p>Hello, World</p>
		`)

		flag := false
		for _, s := range header["Content-Type"] {
			if s == "text/html" {
				flag = true
			}
		}
		assert.Assert(flag)
	})

	t.Run("MatchStaticPath", func(t *testing.T) {
		result, _, err := response("GET", ts.URL, "/static/index.js")
		assert.Eq(err, nil)
		assert.Eq(result, `index.js`)
		result, _, err = response("GET", ts.URL, "/static/css/index.css")
		assert.Eq(err, nil)
		assert.Eq(result, `css/index.css`)
	})

	t.Run("Get", func(t *testing.T) {
		result, _, err := response("GET", ts.URL, "/hello/Danny")
		assert.Eq(err, nil)
		assert.Eq(result, "Hello, Danny")
	})

	t.Run("NoParamNoRoute", func(t *testing.T) {
		result, _, err := response("GET", ts.URL, "/test")
		assert.Eq(err, nil)
		assert.Eq(result, "no param no route")
	})

	t.Run("Post", func(t *testing.T) {
		var jsonStr = []byte(`{"value":"for post"}`)
		req, err := http.NewRequest("POST", ts.URL+"/test/post", bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		var response ForPost
		err = json.NewDecoder(resp.Body).Decode(&response)

		assert.Eq(err, nil)
		assert.Eq(resp.Header.Get("Content-Type"), "application/json")
		assert.Eq(response.Val, "response")
	})

	t.Run("Patch", func(t *testing.T) {
		result, _, err := request("PATCH", ts.URL, "/test/patch", url.Values{
			"value": {"patch"},
		})
		assert.Eq(err, nil)
		assert.Eq(result, "patch")
	})

	t.Run("Query", func(t *testing.T) {
		result, _, err := response("GET", ts.URL, "/test/query?name=Danny")
		assert.Eq(err, nil)
		assert.Eq(result, "Danny")
	})

	t.Run("Handle404NotFound", func(t *testing.T) {
		result, _, err := response("GET", ts.URL, "/404")
		assert.Eq(err, nil)
		assert.Eq(result, "<h1>Page Not Found</h1>")
	})
}

func post(serverUrl, route string, values url.Values) (string, http.Header, error) {
	resp, err := http.PostForm(serverUrl+route, values)
	if err != nil {
		return "", http.Header{}, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", http.Header{}, err
	}
	err = resp.Body.Close()
	if err != nil {
		return "", http.Header{}, err
	}

	return string(b), resp.Header, nil
}

func request(method, serverUrl, route string, values url.Values) (string, http.Header, error) {
	body := strings.NewReader(values.Encode())
	req, err := http.NewRequest(method, serverUrl+route, body)
	if err != nil {
		return "", http.Header{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return "", http.Header{}, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", http.Header{}, err
	}
	err = resp.Body.Close()
	if err != nil {
		return "", http.Header{}, err
	}
	return string(b), resp.Header, nil
}

func response(method, serverUrl, route string) (string, http.Header, error) {
	req, err := http.NewRequest(method, serverUrl+route, bytes.NewBuffer([]byte(``)))
	if err != nil {
		return "", http.Header{}, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", http.Header{}, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", http.Header{}, err
	}
	err = resp.Body.Close()
	if err != nil {
		return "", http.Header{}, err
	}

	return string(b), resp.Header, nil
}
