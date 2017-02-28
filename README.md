# Introduction

looli is a minimalist web framework for go

# Installation

```sh
go get github.com/cssivision/looli
```

# Usage

## Router

looli build on the top of [router](https://github.com/cssivision/router) library, which support `Named parameters` `Wildcard parameters` `Trailing slash redirect` `Case sensitive` `Prefix router` for [detail](https://github.com/cssivision/router).

### Using GET, POST, PUT, PATCH, DELETE and OPTIONS

```go
package main

import (
	"github.com/cssivision/looli"
	"log"
)

func main() {
	router := looli.Default()

	router.Get("/someGet", func(c *looli.Context) {})
	router.Post("/somePost", func(c *looli.Context) {})
	router.Put("/somePut", func(c *looli.Context) {})
	router.Delete("/someDelete", func(c *looli.Context) {})
	router.Patch("/somePatch", func(c *looli.Context) {})
	router.Head("/someHead", func(c *looli.Context) {})
	router.Options("/someOptions", func(c *looli.Context) {})

	log.Fatal(router.Run(":8080"))
}

```
