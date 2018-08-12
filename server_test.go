package rocket_test

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dannypsnl/rocket"

	"github.com/dannypsnl/assert"
)

type User struct {
	Name string `route:"name"`
}

var helloName = rocket.Get("/:name", func(u *User) string {
	return "Hello, " + u.Name
})

func TestServer(t *testing.T) {
	assert := assert.NewTester(t)

	rk := rocket.Ignite(":8080").
		Mount("/hello", helloName)
	ts := httptest.NewServer(rk)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/hello/Danny")
	if err != nil {
		log.Fatal(err)
	}
	greeting, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	assert.Eq(string(greeting), "Hello, Danny")

}
