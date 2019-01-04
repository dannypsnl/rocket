package rocket_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dannypsnl/rocket"
	"github.com/dannypsnl/rocket/fairing"

	"github.com/dannypsnl/assert"
	"github.com/gavv/httpexpect"
)

var (
	forTestHandler = rocket.Get("/", func() string { return "" })
)

func TestOptionsMethod(t *testing.T) {
	rk := rocket.Ignite(":8081").
		Mount("/", forTestHandler)
	ts := httptest.NewServer(rk)
	defer ts.Close()
	e := httpexpect.New(t, ts.URL)

	e.OPTIONS("/").
		Expect().
		Header("Allow").
		Equal("OPTIONS, GET")
}

type Recorder struct {
	fairing.Fairing

	RecordRequestURL []string
}

func (r *Recorder) OnRequest(req *http.Request) *http.Request {
	r.RecordRequestURL = append(r.RecordRequestURL, req.URL.String())
	return req
}

func TestRecorder(t *testing.T) {
	assert := assert.NewTester(t)

	recorder := &Recorder{
		RecordRequestURL: make([]string, 0),
	}

	rk := rocket.Ignite(":9090").
		Attach(recorder).
		Mount("/", rocket.Get("/", func() string { return "home" }))

	ts := httptest.NewServer(rk)
	defer ts.Close()
	e := httpexpect.New(t, ts.URL)

	e.GET("/").
		Expect().Status(http.StatusOK)

	assert.Eq(recorder.RecordRequestURL[0], "/")
}
