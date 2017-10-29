package rocket

import "testing"

func TestHttpMethodHandlerCreaterAPI(t *testing.T) {
	hello := Get("/", func(Context) Response { return "hello" })
	if hello.method != "GET" {
		t.Error(`function[Get] should create a GET handler, but it's`, hello.method, `handler`)
	}
	hello = Post("/", func(Context) Response { return "hello" })
	if hello.method != "POST" {
		t.Error(`function[Post] should create a POST handler, but it's`, hello.method, `handler`)
	}
}
