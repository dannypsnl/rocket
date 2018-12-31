package fairing_test

import (
	"net/http"
	"testing"

	"github.com/dannypsnl/rocket/fairing"
	"github.com/dannypsnl/rocket/response"

	asserter "github.com/dannypsnl/assert"
)

type fakeFairing struct {
	fairing.Fairing
}

func TestFairingAssign(t *testing.T) {
	var f fairing.FairingInterface
	fakeF := &fakeFairing{}
	f = fakeF
	_ = f
}

type counter struct {
	fairing.Fairing

	count uint64
}

func (c *counter) OnRequest(req *http.Request) *http.Request {
	c.count++
	return req
}

func TestFairingCounter(t *testing.T) {
	assert := asserter.NewTester(t)
	var c fairing.FairingInterface = &counter{count: 0}
	c.OnRequest(&http.Request{})
	assert.Eq(c.(*counter).count, uint64(1))
}

func TestDefaultFairing(t *testing.T) {
	assert := asserter.NewTester(t)

	var f fairing.FairingInterface = fairing.Fairing{}
	request := &http.Request{}
	requestViaDefaultFairing := f.OnRequest(request)
	assert.Eq(request, requestViaDefaultFairing)

	resp := response.New("")
	respViaDefaultFairing := f.OnResponse(resp)
	assert.Eq(resp, respViaDefaultFairing)
}
