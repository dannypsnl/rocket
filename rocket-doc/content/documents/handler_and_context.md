---
title: "Request handler & Context"
date: 2018-09-25T21:33:47+08:00
weight: 8
draft: false
---

Rocket's handler contains two parts.

- variant route
- handle function

Basically, we have the creator for handler. It uses like:

```go
package example

import (
	"github.com/dannypsnl/rocket"
)

func handler() string { return "" }

func main() {
	// In `Mount`
    rocket.Get("/hello", handler)
}
```

Now we have a handler function `handler`, and we can use a route `"/hello` and `handler` to create a new Rocket's handler.
When request path matches this route, the response is response of handler function.

We have following creator mapping to HTTP method currently.

- `Get` HTTP Method **GET**
- `Post` HTTP Method **POST**
- `Put` HTTP Method **PUT**
- `Patch` HTTP Method **PATCH**
- `Delete` HTTP Method **DELETE**

Now you already know how to create handler has different method. Let's look the most interesting feature of rocket **User-defined Context**

So the question is, how to have one in rocket?

```go
package example

type User struct {
    Name string `route:"name"`
    Age  uint64 `route:"age"`
}
```

Don't be surprised, that's all, and then we use the type you create as parameter of your handle function

```go
package example

import (
    "github.com/dannypsnl/rocket"
)

func main() {
	// In `Mount`
    rocket.Get("/:name/:age", func (u *User) string {
        return "Hello " + u.Name + ", your age is " + strconv.FormatUint(u.Age, 10)
    })
}
```

Ok, we know how to use the field of context, but where is it came from?

Let's return to `"/:name/:age"`, this is how we fill your context, in variant route, what you defined as `:key` things, will be the value of tag `route:"key"`

At here, we got `route:"name"` & `route:"age"`, so request path `Danny/21` will let your context got string `Danny` & uint64 `21`

p.s. type of `Age` this field is `uint64`, so we will try to parsing the value of request path.
If it's not an `uint64`, then we return **HTTP Status Code: 400**.

But just `route`? Nope, we also have:

- `query:"key"` for request path `/path/to/route?key=value`
- `form:"key"` for FORM request
- `multiform:"key"` for multiple forms request
- `multiform:"key" file:"yes"` for multiple forms request, and it's a file. In this case, the type of field must be `io.ReadCloser`.
- `json:"key"` for request body is JSON(here has some problem, we can also handle GET method just need JSON body,
  this is a bug, it should only work with POST, PUT, PATCH with application/json)
- `header:"key"`, to getting header like `Content-Type`, `Authorization`
- `cookie:"key"`, a very important fact of cookie tag is it only accept you use `*http.Cookie` as field type,
  e.g.

  ```go
  package example

  import "net/http"

  type UserToken struct {
      token *http.Cookie `cookie:"token"`
  }
  ```

  and it would get the most matched cookie, so avoid using duplicate cookie name in your application would be better,
  another important thing is if no cookie matched, this tag won't follow the optional contract,
  so don't use cookie tag in a general purpose context would help you avoid **Bad Request**

#### Optional field

We still have one thing haven't been mentioned, **optional field**, just like its name,
it allowed you omit the field and won't cause **HTTP Status Code: 400**.

Example context definition:

```go
type Transaction struct {
    Amount   uint64 `form:"amount"`
    Canceled *bool  `form:"canceled"`
}
```

At here, the field `Canceled` would be `nil` if you didn't give it a value.

#### Multiple Contexts

Multiple contexts allows you put more than one contexs in your handler function,
this make you can reuse more contexts, for example, here is a proxy of kubernetes List API:

```go
import (
    "net/http"

    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    scheme "k8s.io/client-go/kubernetes/scheme"

    "github.com/dannypsnl/rocket"
    "github.com/dannypsnl/rocket/response"
)

type (
    CheckWatch struct {
        IsWatch *bool `query:"watch"`
    }
    Resource struct {
        Namespace string `route:"namespace"`
        Kind      string `route:"kind"`
    }
)

// ignore all errors
func kubernetesProxy(c *CheckWatch, r *Resource) response.Json {
    opts := &metav1.ListOptions{}
    if c.IsWatch != nil && *c.IsWatch {
        return response.Stream(func(w http.ResponseWriter) {
            opts.Watch = true
            watchInterface, err := kubeClient.CoreV1().RESTClient().Get().
                Namespace(r.Namespace).
                Resource(r.Kind).
                VersionedParams(opts, scheme.ParameterCodec).
                Watch()
            for {
                select {
                case event := <- watchInterface.ResultChan():
                    data, err := json.Marshal(event.Object)
                    w.Write(data)
                }
            }

        }).
            Headers(response.Headers{
                "Content-Type": "application/json",
            })
    }
    // We should create result by Kind actually, but just let me use hard code here as an example
    result := &corev1.EndpointsList{}
    // ignore how to initialize kubeClient
    err := kubeClient.CoreV1().RESTClient().Get().
        Namespace(r.Namespace).
        Resource(r.Kind).
        VersionedParams(opts, scheme.ParameterCodec).
        Do().
        Into(result)
    data, err := json.Marshal(result)
    return response.Json(data)
}

rocket.
    // ignore
    Mount(
        rocket.Get("/api/v1/namespaces/:namespace/:kind", kubernetesProxy),
    ).
    Launch()
```
