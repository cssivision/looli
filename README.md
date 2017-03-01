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

```go
package main

import (
    "net/http"
    "github.com/cssivision/looli"
)

func main() {
    router := looli.Default()

    router.Get("/a/*name", func(c *looli.Context) {
        c.Status(200)
        c.String("hello " + c.Param("name") + "!\n")
    })

    http.ListenAndServe(":8080", router)
}
```
