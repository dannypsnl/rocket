# rocket

![Build Status](https://github.com/dannypsnl/rocket/workflows/Go/badge.svg?branch=master)
[![Build status](https://ci.appveyor.com/api/projects/status/pftm1me961io7hg4?svg=true)](https://ci.appveyor.com/project/dannypsnl/rocket)
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
package main

import (
    "fmt"

    "github.com/dannypsnl/rocket"
)

type User struct {
    Name string `route:"name"`
    Age int `route:"age"`
}

func hello (u *User) string {
    return fmt.Sprintf(
        "Hello, %s.\nYour age is %d.",
        u.Name, u.Age)
}

func main() {
	rocket.Ignite(":8080").
        Mount(
            rocket.Get("/name/:name/age/:age", hello),
        ).
        Launch()
}
```

- First argument of handler creator is a route string can have variant part.
- Second argument is handler function.
  - handler function can have several parameters, these types are you defined to be request context
    Tag in your type will load request value into it!
    - route tag is `route:"name"`, if route contains `/:name`, then value is request URL at this place
      e.g. `/Danny` will let value of `name` is string `Danny`
    - form tag is `form:"key"`, it gets form value from form request
    - form tag is `multiform:"key" limit:"10"`, it gets multiple forms value from form request, download whole file as string
    - json tag is `json:"key"`, it gets POST/PUT body that is JSON
    - header tag is `header:"key"`, it gets header value by key you given from Header
    - cookie tag is `cookie:"key"`, it gets a `*http.Cookie` name is same as key you provided
      remember this field must be `*http.Cookie`
  - return type of handler function is meaningful
    - `response.Html`: returns text as HTML(set Content-Type to `text/html`)
    - `response.Json`: returns text as JSON(set Content-Type to `application/json`)
    - `string`: returns text as plain text(set Content-Type to `text/plain`)
- handler creator name matchs to HTTP Method

#### Mount and Start

```go
rocket.Ignite(":8080"). // Setting port
    Mount(
        rocket.Get("/", index),
        rocket.Get("/static/*path", static),
        rocket.Get("/hello", hello),
    ).
    Launch() // Start Serve
```

- func Ignite get a string to describe port.
- Launch start the server.
- **Mount** receive handlers that exactly handle route. You can emit 0 to N handlers in one `Mount`

#### Fairing

- **OnResponse** and **OnRequest**
  An easy approach:

  ```go
  package example

  import (
      "log"
      "net/http"

      "github.com/dannypsnl/rocket"
      "github.com/dannypsnl/rocket/response"
  )

  type Logger struct {
      rocket.Fairing
  }
  func (l *Logger) OnRequest(r *http.Request) *http.Request {
      log.Printf("request: %#v\n", r)
      return r
  }
  func (l *Logger) OnResponse(r *response.Response) *response.Response {
      log.Printf("response: %#v\n", r)
      return r
  }
  func main()  {
      rocket.Ignite(":6060").
          // Use Attach to emit, you can call Attach multiple time, but carefully at modify data, that might cause problem
          Attach(&Logger{}).
          // Mount(...)
          Launch()
  }
  ```

#### Guard

Guard should be implemented by your **UserDefinedContext**.
Here is an easy example:

```go
package main

import (
	"github.com/dannypsnl/rocket"
)

type User struct {
	Auth *string `header:"Authorization"`
}

func (u *User) VerifyRequest() error {
	// Assuming we have a JWT verify helper function
	if verifyAuthByJWT(u.Auth) {
		return nil
	}
	return rocket.AuthError("not allowed")
}

func main() {
	rocket.Ignite(":8080").
		Mount(
			rocket.Get("/user_data", handler).Guard(&User{}),
		).
		Launch()
}
```
