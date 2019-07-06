package rocket_test

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dannypsnl/rocket"

	"github.com/gavv/httpexpect"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/http2"
)

func TestHTTPsAndHTTP2(t *testing.T) {
	rk := rocket.Ignite(":8082").
		Mount(rocket.Get("/", func() string { return "home" }))
	ts := httptest.NewUnstartedServer(rk)
	ts.TLS = &tls.Config{
		CipherSuites: []uint16{tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256},
		NextProtos:   []string{http2.NextProtoTLS},
	}
	ts.StartTLS()
	defer ts.Close()

	// Create a pool with the server certificate since it is not signed
	// by a known CA
	caCert := ts.Certificate()
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert.RawTBSCertificate)

	// Create TLS configuration with the certificate of the server
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		RootCAs:            caCertPool,
	}
	client := &http.Client{}
	client.Transport = &http2.Transport{
		TLSClientConfig: tlsConfig,
	}

	resp, err := client.Get(ts.URL)
	if err != nil {
		t.Fatalf("failed at GET, error: %s", err)
	}
	if resp.StatusCode != 200 {
		t.Error("request fail")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed at read response body, error: %s", err)
	}
	resp.Body.Close()
	if string(body) != "home" {
		t.Error("response body is wrong")
	}
	if resp.Proto != "HTTP/2.0" {
		t.Error("protocol should be HTTP/2")
	}
}

var (
	forTestHandler = rocket.Get("/", func() string { return "" })
)

func TestOptionsMethod(t *testing.T) {
	rk := rocket.Ignite(":8081").
		Mount(forTestHandler)
	ts := httptest.NewServer(rk)
	defer ts.Close()
	e := httpexpect.New(t, ts.URL)

	e.OPTIONS("/").
		Expect().
		Header("Allow").
		Equal("OPTIONS, GET")
}

type Recorder struct {
	rocket.Fairing

	RecordRequestURL []string
}

func (r *Recorder) OnRequest(req *http.Request) *http.Request {
	r.RecordRequestURL = append(r.RecordRequestURL, req.URL.String())
	return req
}

func TestRecorder(t *testing.T) {
	recorder := &Recorder{
		RecordRequestURL: make([]string, 0),
	}

	rk := rocket.Ignite(":9090").
		Attach(recorder).
		Mount(rocket.Get("/", func() string { return "home" }))

	ts := httptest.NewServer(rk)
	defer ts.Close()
	e := httpexpect.New(t, ts.URL)

	e.GET("/").
		Expect().Status(http.StatusOK)

	assert.Equal(t, "/", recorder.RecordRequestURL[0])
}

type AccessCookie struct {
	Token *http.Cookie `cookie:"token"`
}

func TestGetCookieByUserDefinedContext(t *testing.T) {
	rk := rocket.Ignite("").
		Mount(rocket.Get("/", func(cookie *AccessCookie) string {
			if cookie.Token == nil {
				return "token is nil"
			}
			return cookie.Token.Value
		}))

	ts := httptest.NewServer(rk)
	defer ts.Close()
	e := httpexpect.New(t, ts.URL)

	e.GET("/").WithCookie("token", "123456").
		Expect().Status(http.StatusOK).
		Body().Equal("123456")
}

type AccessHeader struct {
	Auth string `header:"Authorization"`
}

func TestGetHeaderByUserDefinedContext(t *testing.T) {
	rk := rocket.Ignite("").
		Mount(rocket.Get("/", func(header *AccessHeader) string {
			return header.Auth
		}))

	ts := httptest.NewServer(rk)
	defer ts.Close()
	e := httpexpect.New(t, ts.URL)

	e.GET("/").WithHeader("Authorization", "Bear jwt.token.lalala").
		Expect().Status(http.StatusOK).
		Body().Equal("Bear jwt.token.lalala")
}
