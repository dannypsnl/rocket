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

## v0.12.9
