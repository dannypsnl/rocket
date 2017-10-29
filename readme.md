# rocket
Rocket is a web framework inspired by [rocket-rs](https://github.com/SergioBenitez/Rocket).
## Install
Use go get.<br>
`go get github.com/dannypsnl/rocket`
## Usage
### example
You can find example at example folder.
#### Import
```go
import (
    rk "github.com/dannypsnl/rocket"
)
```
#### Create Handler
```go
import "fmt"

var hello = rk.Get("/name/:name/age/:age", func(ctx rk.Context) rk.Response {
    return fmt.Sprintf("Hello, %s.\nYour age is %s.", ctx["name"], ctx["age"])
})

var static = rk.Get("/*path", func(ctx rk.Context) rk.Response {
    return "static"
})

var API = rk.Post("/", func(ctx rk.Context) rk.Response {
    return "Something..."
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
