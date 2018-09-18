package rocket_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/dannypsnl/rocket"
)

var rk *rocket.Rocket

func TestMain(m *testing.M) {
	rk = rocket.Ignite(":8080")
	m.Run()
}

func BenchmarkWithoutUserDefinedContext(b *testing.B) {
	rk = rk.Mount("/home", rocket.Get("/", func() string {
		return "welcome"
	}))
	Request(b, rk, "GET", "/home", nil)
}

func BenchmarkWithUserDefinedContext(b *testing.B) {
	type User struct {
		Name string `route:"name"`
	}

	rk = rk.Mount("/hello", rocket.Get("/:name", func(user *User) string {
		return "welcome-" + user.Name
	}))
	Request(b, rk, "GET", "/hello/kw", nil)
}

func BenchmarkWithCustomResponse(b *testing.B) {
	rk = rk.Mount("/home", rocket.Get("/", func() *rocket.Response {
		return rocket.NewResponse(`welcome-custom-response`).Headers(
			rocket.Headers{
				"Access-Control-Allow-Origin": "*",
			},
		)
	}))
	Request(b, rk, "GET", "/home", nil)
}

func BenchmarkWithHeader(b *testing.B) {
	rk = rk.Mount("/home", rocket.Get("/", func(header *rocket.Header) string {
		return "Content-Type-is-" + header.Get("Content-Type")
	}))
	Request(b, rk, "GET", "/home", nil)
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
