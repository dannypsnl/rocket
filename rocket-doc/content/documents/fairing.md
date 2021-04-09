---
title: "Fairing"
date: 2019-02-02T22:39:18+08:00
weight: 10
draft: false
---

What's fairing? It's an abstraction to avoid over-using the middleware.

But we still need some hooks to record some data or modifying the input/output for certain purpose, to keep this ability, we made fairing, this is how it looks like:

```go
package example

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
// in main function or entrypoint
rocket.Ignite(":8080").
    Attach(&Logger{}).
    // Mount...
    Launch()
```

We can see some points, first, we can implement two kinds of fairing callbacks

- `func OnRequest(r *http.Request) *http.Request`

  this would be called before handlers get request

- `func OnResponse(r *response.Response) *response.Response`

  this would be called after handlers done handling

- `func OnLaunch(r *rocket.Rocket)`

  this would be called at launch time and get meta data of rocket

  For example:

  ```go
  package example

  func (c *Configurator) OnLaunch(r *rocket.Rocket) {
      r.MultiFormBodySizeLimit = 20 // 20 MB
  }
  ```

then we can use the fairing implementor by using `Attach` method to emit it. We can call `Attach` several times, but carefully with it since it could modify request and response!

Why embedded `rocket.Fairing`? It would provide default behavior for `OnRequest` and `OnResponse` if you didn't provide one, so it's a good practice to embedded since we could add more fairing methods into it.
