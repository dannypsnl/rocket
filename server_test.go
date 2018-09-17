package rocket_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dannypsnl/rocket"
	"github.com/gavv/httpexpect"
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
	forPost = rocket.Post("/post", func(f *ForPost) rocket.Json {
		return `{"value": "response"}`
	})
	forPatch = rocket.Patch("/patch", func(f *ForPatch) string {
		return f.Val
	})
	user = rocket.Get("/name/:name/", func(u *User) string {
		return u.Name
	})
	query = rocket.Get("/query", func(u *User) string {
		return u.QueryName
	})
	endWithSlash = rocket.Get("/end-with-slash/", func() string {
		return "you found me"
	})
	handleCookies = rocket.Get("cookies", func(cs *rocket.Cookies) string {
		_, err := cs.Cookie("brabrabra")
		if err == nil {
			return "incorrect!"
		}
		return "cookies"
	})
	customResponseHeader = rocket.Get("/", func() *rocket.Response {
		body := rocket.Json(`{"msg": "welcome"}`)
		return rocket.NewResponse(body).Headers(
			rocket.Header{"Access-Control-Allow-Origin", "*"},
		)
	})
	handlerHeaders = rocket.Get("/headers", func(headers *rocket.Headers) string {
		if headers.Get("x-token") == "token" {
			return "received token"
		}
		return "not receive token"
	})
)

func TestServer(t *testing.T) {
	rk := rocket.Ignite(":8080").
		Mount("/", homePage, staticFiles).
		Mount("/users", user).
		Mount("/test", query, endWithSlash, forPatch, forPost, handleCookies, handlerHeaders).
		Mount("/custom-response-header", customResponseHeader).
		Default(func() rocket.Html {
			return "<h1>Page Not Found</h1>"
		})
	ts := httptest.NewServer(rk)
	defer ts.Close()

	e := httpexpect.New(t, ts.URL)

	t.Run("GetUserName", func(t *testing.T) {
		e.GET("/users/name/Danny").
			Expect().Status(http.StatusOK).
			Body().Equal("Danny")
	})

	t.Run("GetHomePage", func(t *testing.T) {
		e.GET("/").
			Expect().Status(http.StatusOK).
			ContentType("text/html", "").
			Body().Equal(`
		<h1>Title</h1>
		<p>Hello, World</p>
		`)
	})

	t.Run("MatchPathParameter", func(t *testing.T) {
		e.GET("/static/index.js").
			Expect().Status(http.StatusOK).
			Body().Equal(`index.js`)
		e.GET("/static/css/index.css").
			Expect().Status(http.StatusOK).
			Body().Equal(`css/index.css`)
	})

	t.Run("Post", func(t *testing.T) {
		jsonObj := map[string]interface{}{
			"value": "for post",
		}
		expected := map[string]interface{}{
			"value": "response",
		}
		e.POST("/test/post").WithJSON(jsonObj).
			Expect().Status(http.StatusOK).
			ContentType("application/json", "").
			JSON().Equal(expected)
	})

	t.Run("Patch", func(t *testing.T) {
		e.PATCH("/test/patch").WithFormField("value", "patch").
			Expect().Status(http.StatusOK).
			Body().Equal("patch")
	})

	t.Run("Query", func(t *testing.T) {
		e.GET("/test/query").WithQuery("name", "Danny").
			Expect().Status(http.StatusOK).
			Body().Equal("Danny")
	})

	t.Run("Cookies", func(t *testing.T) {
		e.GET("/test/cookies").
			Expect().Status(http.StatusOK).
			Body().Equal("cookies")
	})

	t.Run("EndWithSlash", func(t *testing.T) {
		e.GET("/test/end-with-slash").
			Expect().Status(http.StatusOK).
			Body().Equal("you found me")
	})

	t.Run("Handle404NotFound", func(t *testing.T) {
		e.GET("/404").
			Expect().Status(http.StatusNotFound).
			Body().Equal("<h1>Page Not Found</h1>")
	})

	t.Run("customResponseHeader", func(t *testing.T) {
		expected := map[string]interface{}{
			"msg": "welcome",
		}
		e.GET("/custom-response-header").
			Expect().Status(http.StatusOK).
			Header("Access-Control-Allow-Origin").Equal("*")
		e.GET("/custom-response-header").
			Expect().Status(http.StatusOK).
			JSON().Equal(expected)
	})

	t.Run("Header", func(t *testing.T) {
		expected := "received token"
		e.GET("/test/headers").
			WithHeader("x-token", "token").
			Expect().Status(http.StatusOK).
			Body().Equal(expected)
	})
}
