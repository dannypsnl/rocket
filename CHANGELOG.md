## Latest

- fix: wildcard route matching
- `func (*cookie.Cookie) MaxAge(int)`: use to modified max age field of the cookie
- feat: `func File(filepath string) *Response` at package `response`, use to create a file response with default content-type
- (#126) fix: duplicate path would panic now
- (#125) feature: support auto implements OPTIONS method
- (#86) feature: optional field
- (#134) feature: new fairing

	Now fairing is looking like:
	```go
	import "github.com/dannypsnl/rocket/fairing"

	type YourFairing struct {
		fairing.Fairing
	}

	rocket.Ignite(":8080").
		Attach(&YourFairing{})
		// Ignore
	```
	And purpose would be more like logger than guard

- fix: "/" would let handler creating task fail since out of index
- fix: let status code of response can't be changed twice
- (#87) feature: multiple contexts
	```go
	rocket.Get(func(ctx1 *Ctx1, ctx2 *Ctx2) string {
		// ignore...
	})
	```
- (#152) feature: HTTP/1.1 streaming
	```go
	import (
		"net/http"

		"github.com/dannypsnl/rocket/response"
	)

	// In your handler function
	return response.
		Stream(func(w http.ResponseWriter) {
			w.Write("hello again\n")
			w.Write("and again")
		})
	```

## v0.12.9
