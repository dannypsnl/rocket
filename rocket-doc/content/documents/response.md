---
title: "Response"
date: 2018-09-28T23:55:05+08:00
weight: 9
draft: false
---

> Note: In following context, we only show handle function

> Warning: Don't use `*response.Response` cross goroutine, it's not designed to be thread-safe(and it don't have to be)

Basically rocket contains some type to help you return value more easier

For example, `response.Html`

```go
func() response.Html {
    return `<h1>Title</h1>`
}
```

Then your response when with header `Content-Type` is `text/html`

Currently we have:

- `response.Html`, return `text/html`
- `response.Json`, return `application/json`
  ```go
  func() response.Json {
      return `
      {
          "just": "a json"
      }
      `
  }
  ```
- `string`, return `text/plain`

`Response` is defined under `github.com/dannypsnl/rocket/response` this package

```go
import "github.com/dannypsnl/rocket/response"

func() *response.Response {
    html := response.Html(`
        <h1>Title</h1>
    `)
    return response.New(html)
}
```

As you can see, you can keep your awesome response type feature with `Response`

Factory functions:

- `New`, accept a response type of rocket

  ```go
  response.New("what your user get")
  ```

- `Redirect`, redirect to provided path
  ```go
  response.Redirect("/")
  ```
- `File`, create a reponse from file
  ```go
  response.File("/path/to/file")
  ```
- `Stream`, create a streamable responder by allowing you keep writing data into `http.ResponseWriter`

  ```go
  response.Stream(func(w http.ResponseWriter) (keep bool) {
      _, err := w.Write([]byte(`hello\n`))
      if err != nil {
          return false
      }
      return true
  })
  ```

  This is because Go `http` package help you could use HTTP/1.1 connection as streaming by ignoring **EOF**,
  and somehow we found this is really useful so we keep this ability in rocket

Here is all methods of `Response`:

- `Headers`, accept a header map

  ```go
  response.New("").Headers(response.Headers{
      "Access-Control-Allow-Origin": "*",
  })
  ```

- `Cookies`, accept a cookie list

  ```go
  response.New("").Cookies(
      cookie.New("a", "cookie").
              Expires(time.Now().Add(time.Hour * 24)),
      cookie.New("more", "cookie").
              Expires(time.Now().Add(time.Hour * 24)),
  )
  ```

- `Status`, accept a new status code, it would panic when you give a invalid status code(by RFC, it should be a 3 digit integer) or you try to rewrite it

  ```go
  response.New("Bas request").
      Status(http.StatusBadRequest)
  ```

- `ContentType`, let you modify content-type of response easier(compare to setting **header** directly)

  ```go
  response.Stream(func(w http.ResponseWriter) bool {
      // ignore
  }).
      ContentType("application/json")
  ```
