# rocket

[![CircleCI](https://circleci.com/gh/dannypsnl/rocket.svg?style=svg)](https://circleci.com/gh/dannypsnl/rocket)
[![Go Report Card](https://goreportcard.com/badge/github.com/dannypsnl/rocket)](https://goreportcard.com/report/github.com/dannypsnl/rocket)
[![Coverage Status](https://coveralls.io/repos/github/dannypsnl/rocket/badge.svg?branch=master)](https://coveralls.io/github/dannypsnl/rocket?branch=master)
[![GoDoc](https://godoc.org/github.com/dannypsnl/rocket?status.svg)](https://godoc.org/github.com/dannypsnl/rocket)

Rocket is a web framework inspired by [rocket-rs](https://github.com/SergioBenitez/Rocket).

## Install

`go get github.com/dannypsnl/rocket`

## Usage

#### Import

```go
package example

import (
    rk "github.com/dannypsnl/rocket"
)
```

#### Create Handler

```go
package example

import (
	"fmt"
	
	rk "github.com/dannypsnl/rocket"
)

type User struct {
	Name string `route:"name"`
	Age string `route:"age"`
}

var hello = rk.Get("/name/:name/age/:age", func(u *User) string {
    return fmt.Sprintf(
    	"Hello, %s.\nYour age is %s.",
    	u.Name, u.Age)
})
```

- First argument of handler creator function is a suffix for routing.
- context help you get parameters those you interest in request URL.
- Get, Post function match http method.

#### Mount and Start

```go
rocket.Ignite(":8080"). // Setting port
    Mount("/", index).
    Mount("/", static).
    Mount("/hello", hello).
    Launch() // Start Serve
```

- func Ignite get a string to describe port.
- Launch start the server.
- Mount receive a prefix route and a routes.Handler to handle route.

##### Note

- Base route can't put parameters part. That is illegal route.
