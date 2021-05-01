---
title: "Server"
date: 2018-09-25T21:16:07+08:00
weight: 7
draft: false
---

To create a server, we have to start from `Ignite`.

```go
rocket.Ignite(8080)
```

I think it won't be too hard to notice `8080` means listen port `8080`.

Then we use `Mount` to mount some handlers.

```go
rocket.Ignite(8080).
    Mount(handler)
```

The thing you should know is you can mount several handlers at one `Mount` call.
For example:

```go
rocket.Ignite(8080).
    Mount(handler1, handler2) // and below
```

And the important thing is we high recommended you writing like:

```go
rocket.Ignite(8080).
    Mount(
        rocket.Get("/", handlerFunction),
    )
```

To make route visible when you defining them.

Next is handling **Not Found: 404**, we use `Default` to handle this.

```go
rocket.Ignite(8080).
    // some mounts
    Default(func() response.Html {
        return `<h1>Page Not Found</h1>`
    })
```

Then when rocket can't find any route in router, it will use this function's response.
This is optional, so you can omit it, we have default for default, lol.

p.s. `response.Html` is response magic in rocket, it will set header `Content-Type` as `text/html`.
Then you will see the browser render respnose as HTML

Final, we start our server.

```go
rocket.Ignite(8080).
    // some mounts & default
    Launch()
```

Call `Launch` will start our server, now you can use any HTTP client to see `localhost:8080`

#### Additional API

`EnableHTTPs`, we can create a https server by calling `EnableHTTPs`.

```go
rocket.Ignite(8080).
	EnableHTTPs("cert.pem", "key.pem").
	Launch()
```

Parameters are same as `func (srv *Server) ListenAndServeTLS(certFile, keyFile string) error`.
