---
title: "Guard"
date: 2019-02-17T14:09:30+08:00
weight: 11
draft: false
---

A guard is a context that implemented the interface `rocket.Guard`, it rejects a request by returning a `*response.Response`.

Here is an example:

```go
type User struct {
    Authorization *string `header:"Authorization"`
}

func (u *User) VerifyRequest() *response.Response {
    // Assuming we have a JWT verify helper function
    if verifyAuthByJWT(u.Auth) {
        return nil
    }
    return rocket.AuthError("not allowed")
}

var handler = rocket.Get("/user_data", func(_ *User) string {
    // would return data if `VerifyRequest` do not return any errors
})
```

#### Error response

- `rocket.AuthError`: should be returned when you believe it's an Authorization error, it would bring `403 Forbidden`

  ```go
  rocket.AuthError("auth error, error: %s", err)
  ```

- `rocket.ValidateError`: should be returned when you think the request was something wrong, it would return `400 Bad Request`

  ```go
  rocket.ValidateError("auth error, error: %s", err)
  ```
