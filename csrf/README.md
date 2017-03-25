# csrf

csrf middleware for looli, used to precent [CSRF](https://en.wikipedia.org/wiki/Cross-site_request_forgery) attack.

## Usage

```go
package main

import (
    "net/http"
    "github.com/cssivision/looli"
    "github.com/cssivision/looli/csrf"
    "log"
)

func main() {
    router := looli.Default()

    router.Use(csrf.New(csrf.Options{
        Skip: func(ctx *looli.Context) bool {
            if ctx.Method == http.MethodGet {
                return true
            }
            return false
        },
    }))

    router.Get("/", func(ctx *looli.Context) {
        token := csrf.NewToken(ctx)
        ctx.String("csrf %v\n", token)
    })

    router.Post("/", func(ctx *looli.Context) {
        ctx.String("token valid\n")
    })

    log.Println("server start on http://127.0.0.1:8080")
    log.Fatal(http.ListenAndServe(":8080", router))
}
```

## Parameters

Parameters are passed to the middleware the csrf.New method as follow:

```go
type Options struct {
    // using FormKey to get token
    FormKey   string

    // using HeaderKey to get token
    HeaderKey string

    // using Skip func to check whether to skip csrf check
    Skip      func(*looli.Context) bool

    // parameter below used to store secret key in cookie
    MaxAge    int
    Domain    string
    Path      string
    HttpOnly  bool
    Secure    bool
}
```

## Licenses

All source code is licensed under the [MIT License](https://github.com/cssivision/looli/blob/master/LICENSE).
