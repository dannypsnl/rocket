package fairing_test

import (
	"net/http"
	"testing"

	"github.com/dannypsnl/rocket/fairing"

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
