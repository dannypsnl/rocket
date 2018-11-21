# rocket

[![Build Status](https://travis-ci.org/dannypsnl/rocket.svg)](https://travis-ci.org/dannypsnl/rocket)
[![Go Report Card](https://goreportcard.com/badge/github.com/dannypsnl/rocket)](https://goreportcard.com/report/github.com/dannypsnl/rocket)
[![codecov](https://codecov.io/gh/dannypsnl/rocket/branch/master/graph/badge.svg)](https://codecov.io/gh/dannypsnl/rocket)<Paste>
[![GoDoc](https://godoc.org/github.com/dannypsnl/rocket?status.svg)](https://godoc.org/github.com/dannypsnl/rocket)

Rocket is a web framework inspired by [rocket-rs](https://github.com/SergioBenitez/Rocket).

Document: [https://dannypsnl.github.io/rocket](https://dannypsnl.github.io/rocket)

## Install

`go get github.com/dannypsnl/rocket`

## Usage

#### Import

```go
package example

import (
	"github.com/dannypsnl/rocket"
)
```

#### Create Handler

```go
package example

import (
	"fmt"

	"github.com/dannypsnl/rocket"
)

type User struct {
	Name string `route:"name"`
	Age int `route:"age"`
}

var hello = rocket.Get("/name/:name/age/:age", func(u *User) string {
	return fmt.Sprintf(
		"Hello, %s.\nYour age is %d.",
		u.Name, u.Age)
})
```

- First argument of handler creator is a route string can have variant part. 
- Second argument is handler function.
	- handler function can have a argument, that is a type you define to be request context
		Tag in your type will load request value into it!
		- route tag is `route:"name"`, if route contains `/:name`, then value is request URL at this place
			e.g. `/Danny` will let value of `name` is string `Danny`
		- form tag is `form:"key"`, it get form value from form request
		- json tag is `json:"key"`, it get POST/PUT body that is JSON
	- return type of handler function is meaningful
		- `response.Html`: returns text as HTML(set Content-Type to `text/html`)
		- `response.Json`: returns text as JSON(set Content-Type to `application/json`)
		- `string`: returns text as plain text(set Content-Type to `text/plain`)
- handler creator name is match to HTTP Method

#### Mount and Start

```go
rocket.Ignite(":8080"). // Setting port
	Mount("/", index, static).
	Mount("/hello", hello).
	Launch() // Start Serve
```

- func Ignite get a string to describe port.
- Launch start the server.
- **Mount** receive a base route and a handler that exactly handle route. You can emit 1 to N handlers in one `Mount`
	**Note**: Base route can't put parameters part. That is illegal route.

#### Fairing(experimental release)

- **OnResponse** can help you modified all response at final
	Example:
	```go
	import "github.com/dannypsnl/rocket/response"
	import "github.com/dannypsnl/rocket/fairing"

	rocket.Ignite(":6060").
		Attach(fairing.OnResponse(func(resp *response.Response) *response.Response {
			return resp.Headers(response.Headers{
				"x-fairing": "response",
			})
		})).
		// Mount(...)
		Launch()
	```
