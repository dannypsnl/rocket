package rocket_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dannypsnl/rocket"
	"github.com/dannypsnl/rocket/cookie"
	"github.com/dannypsnl/rocket/fairing"
	"github.com/dannypsnl/rocket/response"
	"github.com/gavv/httpexpect"
)

type (
	User struct {
		Name string `route:"name"`
		Age  int    `route:"age"`

		QueryName string `query:"name"`
	}

	ForPost struct {
		Val string `json:"value"`
	}

	ForPatch struct {
		Val string `form:"value"`
	}

	Files struct {
		FileName string `route:"filename"`
	}

	FilesAndRoute struct {
		V        string `route:"var"`
		FileName string `route:"filename"`
	}
)

type RouteWithJSON struct {
	Field  string `route:"field"`
	Query  string `query:"query_field"`
	JField string `json:"json_field"`
}

var (
	homePage = rocket.Get("/", func() response.Html {
		return `
		<h1>Title</h1>
		<p>Hello, World</p>
		`
	})
	filesAndRoute = rocket.Get("/file/:var/*filename", func(fs *FilesAndRoute) string {
		return fs.V + "/" + fs.FileName
	})
	staticFiles = rocket.Get("/static/*filename", func(fs *Files) string {
		return fs.FileName
	})
	routeWithJSON = rocket.Get("route_with_json/:field", func(r *RouteWithJSON) string {
		return r.Field + r.Query + r.JField
	})
	forPost = rocket.Post("/post", func(f *ForPost) response.Json {
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
		_, err := cs.Get("cookie")
		if err != nil || len(cs.List()) != 1 {
			return "incorrect!"
		}
		return "cookies"
	})
	createCookie = rocket.Get("/new_cookie", func(cs *rocket.Cookies) *response.Response {
		return response.New(``).Cookies(
			cookie.New("set", "set").
				Expires(time.Now().Add(time.Hour * 24)),
		)
	})
	deleteCookie = rocket.Delete("/cookies", func() *response.Response {
		return response.New(``).Cookies(
			cookie.Forget("set"),
		)
	})
	customResponseForHeader = rocket.Get("/", func() *response.Response {
		body := response.Json(`{"msg": "welcome"}`)
		return response.New(body).Headers(
			response.Headers{
				"Access-Control-Allow-Origin": "*",
			},
		)
	})
	handlerHeaders = rocket.Get("/headers", func(header *rocket.Headers) string {
		if header.Get("x-token") == "token" {
			return "received token"
		}
		return "not receive token"
	})
	context = rocket.Get("/context", func(header *rocket.Headers, cookies *rocket.Cookies) string {
		return ""
	})
)

func TestServer(t *testing.T) {
	rk := rocket.Ignite(":8080").
		Mount("/", homePage, staticFiles).
		Mount("/users", user).
		Mount("/test",
			query,
			endWithSlash,
			forPatch,
			forPost,
			handleCookies,
			handlerHeaders,
			context,
			createCookie,
			deleteCookie,
			routeWithJSON,
			filesAndRoute,
		).
		Mount("/custom-response-header", customResponseForHeader).
		Attach(fairing.OnResponse(func(resp *response.Response) *response.Response {
			return resp.Headers(response.Headers{
				"x-fairing": "response",
			})
		})).
		Default(func() response.Html {
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

	t.Run("FilesAndRoute", func(t *testing.T) {
		e.GET("/test/file/css/css/index.css").
			Expect().Status(http.StatusOK).
			Body().Equal("css/css/index.css")
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
		e.GET("/test/cookies").WithCookie("cookie", "cookie").
			Expect().Status(http.StatusOK).
			Body().Equal("cookies")
	})
	t.Run("DeleteCookie", func(t *testing.T) {
		c := e.DELETE("/test/cookies").WithCookie("set", "set").
			Expect().Status(http.StatusOK).
			Cookie("set")

		c.Expires().Equal(time.Unix(0, 0))
	})
	t.Run("CreateNewCookie", func(t *testing.T) {
		startTime := time.Now()

		c := e.GET("/test/new_cookie").
			Expect().Status(http.StatusOK).
			Cookie("set")

		c.Name().Equal("set")
		c.Value()
		c.Path()
		c.Domain()
		c.Expires().InRange(startTime, startTime.Add(time.Hour*24))
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

	t.Run("customResponseForHeader", func(t *testing.T) {
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

	t.Run("Context", func(t *testing.T) {
		e.GET("/test/context").
			Expect().Status(http.StatusOK)
	})

	t.Run("PostHomePage", func(t *testing.T) {
		e.POST("/").
			Expect().Status(http.StatusMethodNotAllowed)
	})

	t.Run("AddHeaderAtResponseFairing", func(t *testing.T) {
		e.GET("/").
			Expect().Status(http.StatusOK).
			Header("x-fairing").Equal("response")
	})

	t.Run("JSONTagCanWorkWithOtherTag", func(t *testing.T) {
		jsonObj := map[string]interface{}{
			"json_field": " json",
		}
		e.GET("/test/route_with_json/field").
			WithQuery("query_field", " query").
			WithJSON(jsonObj).
			Expect().Status(http.StatusOK).
			Body().Equal("field query json")
	})
}
