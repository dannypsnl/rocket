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
			w.Write([]byte("HI\n"))
			w.Write([]byte("Hello again\n"))
			w.Write([]byte("and again\n"))
		})
	```
- (#154) feature: `func (*response.Response) ContentType(contentType string) *response.Response`
- (#129) remove: `rocket.Header`
- (#129) remove: `rocket.Cookies`
- (#129) feature: use User Defined Context to access cookie and header
	```go
	import "net/http"
	// This is your context
	type Ctx struct {
		Auth string `header:"Authorization"`
		// Important thing is cookie only allowed type `*http.Cookie` as field
		// This is because we want to reducing the complex if we introduce
		//	`cookie>value:"token"`
		//	`cookie>expire:"token"`
		// to access cookie sub info
		Token *http.Cookie `cookie:"token"`
	}
	```
- (#147) fix: matching won't fallback bug
- (#116) feature: request guard

## v0.12.9
