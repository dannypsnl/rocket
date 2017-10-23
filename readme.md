# rocket
Rocket is a web framework inspired by [rocket-rs](https://github.com/SergioBenitez/Rocket).
## Install
Use go get.
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

var hello = &rk.Handler {
    Route: "/:name/age/:age",
    Do:    func(ctx rocket.Context) string {
        return fmt.Sprintf("Hello, %s\nYour age is %s", ctx["name"], ctx["age"])
    },
}

var index = rk.Get("/", func(ctx rk.Context) rk.Response {
    return "index"
})
```
- Handler.Route is a suffix for routing.
- context help you get parameters those you interest in request URL.
#### Mount and Start
```go
rocket.Ignite(":8080"). // Setting port
    Mount("/", index).
    Mount("/*path", static).
    Mount("/hello", hello).
    Launch() // Start Serve
```
- func Ignite get a string to describe port.
- Launch start the server.
- Mount receive a prefix route and a routes.Handler to handle route.
