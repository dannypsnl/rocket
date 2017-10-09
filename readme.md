# rocket
This pkg is a web framework.
## Install
Use your go get.
## Usage
#### Import
```go
import (
    "rocket"
    "rocket/routes"
)
```
#### Create Handler
```go
const hello = routes.Handler {
    Route: "",
    Do:    func(...interface{}) string {
        return "Hello!!!"
    },
}
```
- Handler.Route is a suffix for routing.
#### Mount and Start
```go
func index(w http.ResponseWriter, r *http.Request) {
    // ...
}

rocket.Ignite(":8080").
    Mount("/hello", hello).
    MountNative("/", index)
    Launch() // Start Serve
```
- func Ignite get a string to describe port.
- Launch start the server.
- Mount receive a prefix route and a routes.Handler to handle route.
- And you also can use MountNative to keep using native handler.
