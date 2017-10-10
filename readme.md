# rocket
A web framework.
## Install
Use your go get.
## Usage
#### Import
```go
import (
    "rocket"
)
```
#### Create Handler
```go
import "fmt"

var hello = rocket.Handler {
    Route: "/:name/:age",
    Do:    func(Context map[string]string) string {
        return fmt.Sprintf("Hello, %s\nYour age is %s", Context["name"], Context["age"])
    },
}
```
- Handler.Route is a suffix for routing.
#### Mount and Start
```go
rocket.Ignite(":8080").
    Mount("/hello", hello).
    Launch() // Start Serve
```
- func Ignite get a string to describe port.
- Launch start the server.
- Mount receive a prefix route and a routes.Handler to handle route.
