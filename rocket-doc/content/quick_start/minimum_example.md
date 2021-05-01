---
title: "Minimum Example"
date: 2018-09-23T23:51:06+08:00
---

Before you write down any code. You need to import the package.

```go
import (
    "github.com/dannypsnl/rocket"
)
```

With Rocket, you will create a lots of handler function,
here is a basic handler function with user-defined context.

```go
type User struct {
    Name string `route:"name"`
    Age  uint64 `route:"age"`
}
func hello(u *User) string {
    return "Hello " + u.Name + ", your age is " + strconv.FormatUint(u.Age, 10)
}
```

How to let it work?

```go
// main.go
func main() {
    rocket.Ignite(8080).
        Mount(
            // put `hello` under a path `/user/:name/:age`, where `:name` and `:age` are variant parameters
            rocket.Get("/user/:name/:age", hello),
        ).
        Launch()
}
```

Now execute `go run main.go`, open your browser to `localhost:8080/user/Danny/21`.

Then you will see `Hello Danny, your age is 21`.

Or use `curl`:

```bash
$ curl localhost:8080/user/Danny/21
Hello Danny, your age is 21
```
