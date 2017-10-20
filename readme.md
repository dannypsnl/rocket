# rocket
Rocket is a web framework inspired by [rocket-rs](https://github.com/SergioBenitez/Rocket).
## Install
Use go get.
`go get github.com/dannypsnl/rocket`
## Usage
#### Import
```go
import (
    "github.com/dannypsnl/rocket"
)
```
#### Create Handler
```go
import "fmt"

var hello = rocket.Handler {
    Route: "/:name/:age",
    Do:    func(context map[string]string) string {
        return fmt.Sprintf("Hello, %s\nYour age is %s", context["name"], context["age"])
    },
}
```
- Handler.Route is a suffix for routing.
- context help you get parameters those you interest in request URL.
#### Mount and Start
```go
rocket.Ignite(":8080").
    Mount("/", index).
    Mount("/*path", static).
    Mount("/hello", hello).
    Launch() // Start Serve
```
- func Ignite get a string to describe port.
- Launch start the server.
- Mount receive a prefix route and a routes.Handler to handle route.
