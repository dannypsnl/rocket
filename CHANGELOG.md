## Latest

- provide constants in response for content type hint: `"github.com/dannypsnl/rocket/response/content_type"`
    - `content_type.HTML`
    - `content_type.JSON`
    - `content_type.TextPlain`

## v1.1.0

- (#210) Let guard be a context implements `rocket.Guard`
- (#204) `OnClose` takes a function f and runs it after server closed
- (#205) use `int` as **port**, not `string`
- API rename
  - `fairingInterface` -> `Fairing`
  - `Fairing` -> `DefaultFairing`
- (#198) response: redirect
  ```go
  response.Redirect("/")
  ```
- (#194) supports multiple forms

  ```go
  type RequestContext struct {
      Normal string `multiform:"name"`
      File io.ReadCloser `multiform:"file" file:"yes"`
  }
  ```

## v1.0.0

- (#172) improve MIME type detection

## v0.14.0

- (#165) add directly access to `*http.Request`
  ```go
  type RequestContext struct {
  	Request *http.Request `http:"request"`
  }
  rocket.Get("/", func(c *RequestContext) string { return c.Request.URL.Path })
  ```
- add `EnableHTTPs` method for creating HTTPs allowed server
  ```go
  rocket.Ignite(":443").
  	EnableHTTPs("cert.pem", "key.pem").
  	Launch()
  ```
- (#183) fix: reject void function as handler
- add `OnLaunch` fairing

  ```go
  type YourFairing struct {
  	rocket.Fairing
  }

  func (f *YourFairing) OnLaunch(r *rocket.Rocket) {
  	// get rocket structure info at launch time
  }
  ```

- remove subpackage `fairing`
  `fairing.Fairing` ~> `rocket.Fairing`
- remove base route from design, NOTE: it's a big break change

  New style example:

  ```go
  rocket.Ignite(":8080").
  	Mount(
  		rocket.Get("/", home),
  		rocket.Get("/static/*filepath", staticFiles)
   	).
  	Launch()
  ```

## v0.13.0

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
- (#123) feature: HTTP/1.1 streaming

  ```go
  import (
  	"net/http"

  	"github.com/dannypsnl/rocket/response"
  )

  // In your handler function
  return response.
  	Stream(func(w http.ResponseWriter) (keep bool) {
  		// keep writing until failed
  		_, err := w.Write([]byte("HI\n"))
  		if err != nil {
  			return false
  		}
  		return true
  	})
  ```

## v0.12.9
