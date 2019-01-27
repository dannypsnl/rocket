---
title: "Server"
date: 2018-09-25T21:16:07+08:00
weight: 7
draft: false
---

To create a server, we have to start from `Ignite`.
```go
rocket.Ignite(":8080")
```

We use the same pattern as Go native `http` package.

So I think it won't be to hard to notice `:8080` means listen port `8080`.

Then we will use `Mount` mounts some handlers.
```go
rocket.Ignite(":8080").
    Mount("/", handler)
```

The first argument of `Mount` is base route. This is a leading route for following handlers mounted under this route.

For example, if base route is `"/base"`, route of handler is `"/hello"`, the final route is `"/base/hello"`

Then one thing you should know is you can mount several handlers at one `Mount` call.
For example:
```go
rocket.Ignite(":8080").
    Mount("/", handler1, handler2) // and below
```

Next is handling **Not Found: 404**, we use `Default` to handle this.
```go
rocket.Ignite(":8080").
    // some mounts
    Default(func() rocket.Html {
        return `<h1>Page Not Found</h1>`
    })
```

Then when rocket can't find any route in router, it will use this function's response.
This is optional, so you can omit it, we have default for default, lol.

p.s. `rocket.Html` is response magic in rocket, it will set header `Content-Type` as `text/html`.
Then you will see the browser render respnose as HTML

Final, we start our server.
```go
rocket.Ignite(":8080").
    // some mounts & default
    Launch()
```

Call `Launch` will start our server, now you can use any HTTP client to see `localhost:8080`