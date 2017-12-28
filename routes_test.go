package rocket

import "testing"

func TestHttpMethodHandlerCreaterAPI(t *testing.T) {
	hello := Get("/", func(Context) Res { return "hello" })
	if hello.method != "GET" {
		t.Error(`function[Get] should create a GET handler, but it's`, hello.method, `handler`)
	}
	hello = Post("/", func(Context) Res { return "hello" })
	if hello.method != "POST" {
		t.Error(`function[Post] should create a POST handler, but it's`, hello.method, `handler`)
	}
	hello = Put("/", func(Context) Res { return "hello" })
	if hello.method != "PUT" {
		t.Error(`function[Put] should create a PUT handler, but it's`, hello.method, `handler`)
	}
	hello = Delete("/", func(Context) Res { return "hello" })
	if hello.method != "DELETE" {
		t.Error(`function[Delete] should create a DELETE handler, but it's`, hello.method, `handler`)
	}
}
