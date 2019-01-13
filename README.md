# rocket

[![Build Status](https://travis-ci.org/dannypsnl/rocket.svg)](https://travis-ci.org/dannypsnl/rocket)
[![Build status](https://ci.appveyor.com/api/projects/status/pftm1me961io7hg4?svg=true)](https://ci.appveyor.com/project/dannypsnl/rocket)
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
	- handler function can have several parameters, these types are you defined to be request context
		Tag in your type will load request value into it!
		- route tag is `route:"name"`, if route contains `/:name`, then value is request URL at this place
			e.g. `/Danny` will let value of `name` is string `Danny`
		- form tag is `form:"key"`, it get form value from form request
		- json tag is `json:"key"`, it get POST/PUT body that is JSON
		- header tag is `header:"key"`, it can get header value by key you given from Header
		- cookie tag is `cookie:"key"`, it would get a `*http.Cookie` name is same as key you provided
			remember this field must be `*http.Cookie`
	- return type of handler function is meaningful
		- `response.Html`: returns text as HTML(set Content-Type to `text/html`)
		- `response.Json`: returns text as JSON(set Content-Type to `application/json`)
		- `string`: returns text as plain text(set Content-Type to `text/plain`)
- handler creator name matchs to HTTP Method

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

- **OnResponse** and **OnRequest**
	An easy approach:
	```go
	import (
		"net/http"

		"github.com/dannypsnl/rocket/response"
		"github.com/dannypsnl/rocket/fairing"
	)

	type Logger struct {
		fairing.Fairing
	}
	func (l *Logger) OnRequest(r *http.Request) *http.Request {
		log.Printf("request: %#v\n", r)
		return r
	}
	func (l *Logger) OnResponse(r *response.Response) *response.Response {
		log.Printf("response: %#v\n", r)
		return r
	}
	rocket.Ignite(":6060").
		// Use Attach to emit, you can call Attach multiple time, but carefully at modify data, that might cause problem
		Attach(&Logger{}).
		// Mount(...)
		Launch()
	```

#### Guard

Guard should be implemented by your **UserDefinedContext**.
Here is an easy example:
```go
import (
	"errors"
	"net/http"

	"github.com/dannypsnl/rocket"
)

type User struct {}

func (u *User) VerifyRequest(req *http.Request) (rocket.Action, error) {
	user, password, ok := req.BasicAuth()
	if ok && user == "danny" && password == "password" {
		return rocket.Success, nil
	}
	return rocket.Failure, errors.New("not allowed")
}
```
