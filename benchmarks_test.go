package rocket_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/dannypsnl/rocket"
	"github.com/dannypsnl/rocket/response"
)

func BenchmarkRequest(b *testing.B) {
	b.Run("WithoutUserDefinedContext", func(b *testing.B) {
		rk := rocket.Ignite(":8080").
			Mount("/home", rocket.Get("/", func() string {
				return "welcome"
			}))
		Request(b, rk, "GET", "/home", nil)
	})
	b.Run("WithUserDefinedContext", func(b *testing.B) {
		type User struct {
			Name string `route:"name"`
		}

		rk := rocket.Ignite(":8080").
			Mount("/hello", rocket.Get("/:name", func(user *User) string {
				return "welcome-" + user.Name
			}))
		Request(b, rk, "GET", "/hello/kw", nil)
	})
	b.Run("WithCustomResponse", func(b *testing.B) {
		rk := rocket.Ignite(":8080").
			Mount("/home", rocket.Get("/", func() *response.Response {
				return response.New(`welcome-custom-response`).Headers(
					response.Headers{
						"Access-Control-Allow-Origin": "*",
					},
				)
			}))
		Request(b, rk, "GET", "/home", nil)
	})
	b.Run("WithHeader", func(b *testing.B) {
		rk := rocket.Ignite(":8080").
			Mount("/home", rocket.Get("/", func(header *rocket.Headers) string {
				return "Content-Type-is-" + header.Get("Content-Type")
			}))
		Request(b, rk, "GET", "/home", nil)
	})
}

func Request(b *testing.B, rk *rocket.Rocket, method, path string, body io.Reader) {
	b.Helper()
	req, _ := http.NewRequest(method, path, body)
	b.ReportAllocs()
	b.ResetTimer()
	w := &mockRespWriter{}
	for i := 0; i < b.N; i++ {
		rk.ServeHTTP(w, req)
	}
}

type mockRespWriter struct{}

func (m *mockRespWriter) Header() (h http.Header) {
	return http.Header{}
}

func (m *mockRespWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (m *mockRespWriter) WriteHeader(int) {}
