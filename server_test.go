package rocket_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dannypsnl/rocket"

	"github.com/dannypsnl/assert"
)

type User struct {
	Name string `route:"name"`
}

var (
	noParamNoRoute = rocket.Get("/", func() string {
		return "no param no route"
	})
	helloName = rocket.Get("/:name", func(u *User) string {
		return "Hello, " + u.Name
	})
	forPost = rocket.Post("/post", func() string {
		return "for post"
	})
)

func TestServer(t *testing.T) {
	assert := assert.NewTester(t)

	rk := rocket.Ignite(":8080").
		Mount("/", noParamNoRoute).
		Mount("/for", forPost).
		Mount("/hello", helloName).
		Mount("/test", noParamNoRoute)
	ts := httptest.NewServer(rk)
	defer ts.Close()

	t.Run("Get", func(t *testing.T) {
		result, err := response("GET", ts.URL, "/hello/Danny")
		assert.Eq(err, nil)
		assert.Eq(result, "Hello, Danny")
	})

	t.Run("NoParamNoRoute", func(t *testing.T) {
		result, err := response("GET", ts.URL, "/test")
		assert.Eq(err, nil)
		assert.Eq(result, "no param no route")

		result, err = response("GET", ts.URL, "/")
		assert.Eq(err, nil)
		assert.Eq(result, "no param no route")
	})

	t.Run("Post", func(t *testing.T) {
		result, err := response("POST", ts.URL, "/for/post")
		assert.Eq(err, nil)
		assert.Eq(result, "for post")
	})
}

func response(method, serverUrl, route string) (string, error) {
	req, err := http.NewRequest(method, serverUrl+route, bytes.NewBuffer([]byte(``)))
	if err != nil {
		return "", err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	err = resp.Body.Close()
	if err != nil {
		return "", err
	}

	return string(b), nil
}
