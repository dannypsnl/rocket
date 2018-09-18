package rocket_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/dannypsnl/rocket"
)

func BenchmarkWithoutUserDefindContext(b *testing.B) {
	rk := rocket.Ignite(":8080").
		Mount("/home", rocket.Get("/", func() string {
			return "welcome"
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
