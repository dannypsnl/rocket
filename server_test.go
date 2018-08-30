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
}

type ForPost struct {
	Val string `json:"value"`
}

type ForPatch struct {
	Val string `form:"value"`
}

var (
	homePage = rocket.Get("/", func() rocket.Html {
		return `
		<h1>Title</h1>
		<p>Hello, World</p>
		`
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
)

func TestServer(t *testing.T) {
	assert := assert.NewTester(t)

	rk := rocket.Ignite(":8080").
		Mount("/", homePage, forPost).
		Mount("/for", forPatch).
		Mount("/hello", helloName).
		Mount("/test", noParamNoRoute).
		Default(func() rocket.Html {
			return "<h1>Page Not Found</h1>"
		})
	ts := httptest.NewServer(rk)
	defer ts.Close()

	t.Run("GetHTML", func(t *testing.T) {
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
		req, err := http.NewRequest("POST", ts.URL+"/post", bytes.NewBuffer(jsonStr))
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
		result, _, err := request("PATCH", ts.URL, "/for/patch", url.Values{
			"value": {"for patch"},
		})
		assert.Eq(err, nil)
		assert.Eq(result, "for patch")
	})

	t.Run("404", func(t *testing.T) {
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
