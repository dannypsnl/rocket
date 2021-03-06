package rocket_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dannypsnl/rocket"
	"github.com/dannypsnl/rocket/cookie"
	"github.com/dannypsnl/rocket/response"

	"github.com/gavv/httpexpect"
)

type (
	User struct {
		Name string `route:"name"`
		Age  int    `route:"age"`
	}
	Article struct {
		ID string `query:"article_id"`
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
	OptionalContext struct {
		A *string `query:"a"`
	}
)

type RouteWithJSON struct {
	Field  string `route:"field"`
	Query  string `query:"query_field"`
	JField string `json:"json_field"`
}

func homePage() response.Html {
	return `
		<h1>Title</h1>
		<p>Hello, World</p>
		`
}
func staticFiles(fs *Files) string {
	return fs.FileName
}
func user(u *User) string {
	return u.Name
}
func testQuery(u *Article) string {
	return u.ID
}
func filesAndRoute(fs *FilesAndRoute) string {
	return fs.V + "/" + fs.FileName
}
func routeWithJSON(r *RouteWithJSON) string {
	return r.Field + r.Query + r.JField
}
func testPatch(f *ForPatch) string {
	return f.Val
}
func optionalFieldHandler(optionalContext *OptionalContext) string {
	if optionalContext.A == nil {
		return "a is nil"
	}
	return "a is " + *optionalContext.A
}
func testPost(f *ForPost) response.Json {
	return `{"value": "response"}`
}
func endWithSlash() string {
	return "you found me"
}
func deleteCookie() *response.Response {
	return response.New(``).Cookies(
		cookie.Forget("set"),
	)
}
func customResponseForHeader() *response.Response {
	body := response.Json(`{"msg": "welcome"}`)
	return response.New(body).Headers(
		response.Headers{
			"Access-Control-Allow-Origin": "*",
		},
	)
}

var ()

func TestServer(t *testing.T) {
	rk := rocket.Ignite(8080).
		Mount(
			rocket.Get("/", homePage),
			rocket.Get("/static/*filename", staticFiles),
			rocket.Get("/file/:var/*filename", filesAndRoute),
			rocket.Get("/users/name/:name/", user),
			rocket.Get("/test/query", testQuery),
			rocket.Get("/test/end-with-slash/", endWithSlash),
			rocket.Patch("/test/patch", testPatch),
			rocket.Post("/test/post", testPost),
			rocket.Delete("/test/cookies", deleteCookie),
			rocket.Get("/test/route_with_json/:field", routeWithJSON),
			rocket.Get("/test/file/:var/*filename", filesAndRoute),
			rocket.Get("/test/optional/", optionalFieldHandler),
			rocket.Get("/custom-response-header", customResponseForHeader),
		).
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
		e.GET("/file/css/css/index.css").
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
		e.GET("/test/query").WithQuery("article_id", "123").
			Expect().Status(http.StatusOK).
			Body().Equal("123")
	})

	t.Run("DeleteCookie", func(t *testing.T) {
		c := e.DELETE("/test/cookies").WithCookie("set", "set").
			Expect().Status(http.StatusOK).
			Cookie("set")

		c.Expires().Equal(time.Unix(0, 0))
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

	t.Run("PostHomePage", func(t *testing.T) {
		e.POST("/").
			Expect().Status(http.StatusMethodNotAllowed)
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

	t.Run("OptionalField", func(t *testing.T) {
		e.GET("/test/optional").
			Expect().Status(http.StatusOK).
			Body().Equal("a is nil")

		e.GET("/test/optional").WithQuery("a", "a").
			Expect().Status(http.StatusOK).
			Body().Equal("a is a")
	})

	t.Run("MultipleContext", func(t *testing.T) {
		type User struct {
			ID   int    `route:"user_id"`
			Name string `form:"user_name"`
			Age  int    `query:"user_age"`
		}
		type Article struct {
			ID        int    `route:"article_id"`
			Name      string `form:"article_name"`
			CreatedAt string `query:"article_created_at"`
			AuthorID  int    `route:"user_id"`
		}
		rk := rocket.Ignite(8081).
			Mount(rocket.Post("/users/:user_id/articles/:article_id", func(user *User, article *Article) *response.Response {
				resp := response.New("")
				if user.Age != 18 {
					resp.Status(400)
				}
				if user.ID != 1 {
					resp.Status(400)
				}
				if user.Name != "hi" {
					resp.Status(400)
				}
				if article.AuthorID != 1 {
					resp.Status(400)
				}
				if article.CreatedAt != "1994" {
					resp.Status(400)
				}
				if article.ID != 2 {
					resp.Status(400)
				}
				if article.Name != "hello" {
					resp.Status(400)
				}
				return resp
			}))
		ts := httptest.NewServer(rk)
		defer ts.Close()

		e := httpexpect.New(t, ts.URL)

		e.POST("/users/1/articles/2").
			WithFormField("user_name", "hi").
			WithFormField("article_name", "hello").
			WithQuery("user_age", "18").
			WithQuery("article_created_at", "1994").
			Expect().Status(http.StatusOK)
	})
}
