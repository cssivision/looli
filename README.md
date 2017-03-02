# Introduction

looli is a minimalist web framework for go

# Installation

```sh
go get github.com/cssivision/looli
```

# Usage

## Router

looli build on the top of [router](https://github.com/cssivision/router) library, which support `Named parameters` `Wildcard parameters` `Trailing slash redirect` `Case sensitive` `Prefix router`, for [detail](https://github.com/cssivision/router).

### Using GET, POST, PUT, PATCH, DELETE and OPTIONS

```go
package main

import (
	"github.com/cssivision/looli"
	"log"
)

func main() {
	router := looli.Default()

	router.Get("/a", func(c *looli.Context) {})
	router.Post("/a", func(c *looli.Context) {})
	router.Put("/a", func(c *looli.Context) {})
	router.Delete("/a", func(c *looli.Context) {})
	router.Patch("/a", func(c *looli.Context) {})
	router.Head("/a", func(c *looli.Context) {})
	router.Options("/a", func(c *looli.Context) {})

	log.Fatal(router.Run(":8080"))
}

```

### Named parameter

Named parameters only match a single path segment:
```
Pattern: /user/:name

 /user/gordon              match
 /user/you                 match
 /user/gordon/profile      no match
 /user/                    no match

Pattern: /:user/:name

 /a/gordon                 match
 /b/you                    match
 /user/gordon/profile      no match
 /user/                    no match
```

```go
package main

import (
    "net/http"
    "github.com/cssivision/looli"
)

func main() {
    router := looli.Default()

    router.Get("/a/:name", func(c *looli.Context) {
        c.Status(200)
        c.String("hello " + c.Param("name") + "!\n")
    })

    http.ListenAndServe(":8080", router)
}
```

### Wildcard pattern

Match everything, therefore they must always be at the end of the pattern:

```
Pattern: /src/*filepath

 /src/                     match
 /src/somefile.go          match
 /src/subdir/somefile.go   match
```

```go
package main

import (
    "net/http"
    "github.com/cssivision/looli"
)

func main() {
    router := looli.Default()

    router.Get("/a/*filepath", func(c *looli.Context) {
        c.Status(200)
        c.String("hello " + c.Param("filepath") + "!\n")
    })

    http.ListenAndServe(":8080", router)
}
```

### Trailing slash redirect

By default `TrailingSlashRedirect = true` which means if we register path `/a/b`, we can request with `/a/b/`, conversely also success. redirect will work only in the situation that the request can not found.

```
/a/b -> /a/b/
/a/b/ -> /a/b
```

```go
package main

import (
    "net/http"
    "github.com/cssivision/looli"
)

func main() {
    router := looli.Default()

    // default is true, we can forbidden this behavior by set is to false
    // request with /a/ will get 404
    router.SetTrailingSlashRedirect(false)

    router.Get("/a", func(c *looli.Context) {
        c.Status(200)
        c.String("hello world!\n")
    })

    http.ListenAndServe(":8080", router)
}
```

### Case sensitive

By default `IgnoreCase = false`, which means if we register path `/a/b`, request with `/A/B` will get `404 not found`. if we set `IgnoreCase = true`, request with `/A/B` will success.

```go
package main

import (
    "net/http"
    "github.com/cssivision/looli"
)

func main() {
    router := looli.Default()

    // default is false, we can forbidden this behavior by set is to true
    // request with /A/ will success.
    router.SetIgnoreCase(true)

    router.Get("/a", func(c *looli.Context) {
        c.Status(200)
        c.String("hello world!\n")
    })

    http.ListenAndServe(":8080", router)
}
```

### Prefix router

Group router using prefix

```go
package main

import (
    "net/http"
    "github.com/cssivision/looli"
)

func main() {
    router := looli.Default()

    v1 := router.Prefix("/v1")
    v1.Get("/a", func(c *looli.Context) {
        c.Status(200)
        c.String("hello world version1\n")
    })

    v2 := router.Prefix("/v2")
    v2.Get("/a", func(c *looli.Context) {
        c.Status(200)
        c.String("hello world version2\n")
    })

    router.Get("/a", func(c *looli.Context) {
        c.Status(200)
        c.String("hello world!\n")
    })

    http.ListenAndServe(":8080", router)
}
```
