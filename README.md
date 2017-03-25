# looli
[![Build Status](https://img.shields.io/travis/cssivision/looli.svg?style=flat-square)](https://travis-ci.org/cssivision/looli)
[![Coverage Status](http://img.shields.io/coveralls/cssivision/looli.svg?style=flat-square)](https://coveralls.io/github/cssivision/looli?branch=master)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://github.com/cssivision/looli/blob/master/LICENSE)


looli is a minimalist web framework for golang.

# Feature

* [Router](#router)
    * [full method support](#using-get-post-put-patch-delete-and-options)
    * [named parameter](#named-parameter)
    * [wildcard pattern](#wildcard-pattern)
    * [trailing slash redirect](#trailing-slash-redirect)
    * [case sensitive](#case-sensitive)
    * [serving static files](#serving-static-files)
* [Context](#context)
    * [query and form](#query-and-form)
    * [header and cookie](#header-and-cookie)
    * [data binding](#data-binding)
    * [string xml json rendering](#string-json-rendering)
    * [html rendering](#html-rendering)
* [Middleware](#middleware)
    * [using middleware](#using-middleware)
    * [builtin middlewares](#builtin-middlewares)
    * [custome middleware](#custome-middleware)

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

By default will redirect, which means if we register path `/a/b`, we can request with `/a/b/`, conversely also success. redirect will work only in the situation that the request can not found, if both define path `/a/b` and `/a/b/`, redirect will not work.

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

By default is not case sensitive, which means if we register path `/a/b`, request with `/A/B` will get `404 not found`. if we set `true`, request path with `/A/B` will success.

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

### Serving static files

```go
package main

import (
    "net/http"
    "github.com/cssivision/looli"
)

func main() {
    router := looli.Default()

    // Serve file in the path
    router.StaticFile("/somefile.go", "file/path")

    // Serve files in staic directory
    router.Static("/static", "./static")

    http.ListenAndServe(":8080", router)
}
```

## Context

Context supply some syntactic sugar.

### Query and Form

```go
package main

import (
    "net/http"
    "github.com/cssivision/looli"
)

func main() {
    router := looli.Default()

    router.Get("/query", func(c *looli.Context) {
        id := c.Query("id")
        name := c.DefaultQuery("name", "cssivision")
        c.Status(200)
        c.String("hello %s, %s\n", id, name)
    })

    router.Post("/form", func(c *looli.Context) {
        name := c.DefaultPostForm("name", "somebody")
        age := c.PostForm("age")
        c.Status(200)
        c.JSON(looli.JSON{
            "name": name,
            "age": age,
        })
    })

    http.ListenAndServe(":8080", router)
}
```

query
```sh
curl 'localhost:8080/query?id=1&name=cssivision'
```

form
```sh
curl -d 'age=21&other=haha' 'localhost:8080/form?id=1&name=cssivision'
```

### Header and Cookie

Use method to operate header and cookie

```go
package main

import (
    "fmt"
    "github.com/cssivision/looli"
    "log"
    "net/http"
)

func main() {
    router := looli.Default()

    router.Get("/header", func(c *looli.Context) {
        fmt.Println(c.Header("User-Agent"))
        c.SetHeader("fake-header", "fake")
        c.Status(200)
        c.String("fake header has setted\n")
    })

    router.Get("/cookie", func(c *looli.Context) {
        val, _ := c.Cookie("fake-cookie")
        fmt.Println(val)
        c.SetCookie(&http.Cookie{
            Name: "fake-cookie",
            Value: "fake",
        })
        c.Status(200)
        c.String("fake cookie has setted\n")
    })

    log.Fatal(router.Run(":8080"))
}
```

### Data binding

To bind a request into a type, use data binding, data can from query, post body. currently support binding of JSON, XML and standard form values (x-www-form-urlencoded and multipart/form-data). When using the Bind-method, the binder depending on the Content-Type header.

Note that you need to set the corresponding binding tag on all fields you want to bind. For example, when binding from JSON, set json:"fieldname".

```go
package main

import (
    "fmt"
    "github.com/cssivision/looli"
    "net/http"
)

type Infomation struct {
    Name string`json:"name"`
    Age int`json:"age"`
}

func main() {
    router := looli.Default()

    // curl 'localhost:8080/query?name=cssivision&age=21'
    router.Get("/query", func(c *looli.Context) {
        query := new(Infomation)
        if err := c.Bind(query); err != nil {
            fmt.Println(err)
            return
        }
        fmt.Println(query.Name)
        fmt.Println(query.Age)
        c.Status(200)
        c.JSON(query)
    })

    // curl -d "name=cssivision&age=21" 'localhost:8080/form'
    router.Post("/form", func(c *looli.Context) {
        form := new(Infomation)
        if err := c.Bind(form); err != nil {
            fmt.Println(err)
            return
        }
        fmt.Println(form.Name)
        fmt.Println(form.Age)
        c.Status(200)
        c.JSON(form)
    })

    // curl  -H "Content-Type: application/json" -X POST -d '{"name":"cssivision","age":21}' localhost:8080/json
    router.Post("/json", func(c *looli.Context) {
        json := new(Infomation)
        if err := c.Bind(json); err != nil {
            fmt.Println(err)
            return
        }
        fmt.Println(json.Name)
        fmt.Println(json.Age)
        c.Status(200)
        c.JSON(json)
    })

    http.ListenAndServe(":8080", router)
}
```

### String JSON rendering

```go
package main

import (
    "github.com/cssivision/looli"
    "net/http"
)

func main() {
    router := looli.Default()

    router.Get("/string", func(c *looli.Context) {
        c.String("the response is %s\n", "string")
    })

    router.Get("/json1", func(c *looli.Context) {
        c.JSON(looli.JSON{
            "name": "cssivision",
            "age": 21,
        })
    })

    router.Get("/json2", func(c *looli.Context) {
        var msg struct {
            Name string`json:"name"`
            Age int`json:"age"`
        }

        msg.Name = "cssivision"
        msg.Age = 21

        c.JSON(msg)
    })

    http.ListenAndServe(":8080", router)
}
```

### HTML rendering

```go
package main

import (
    "github.com/cssivision/looli"
    "net/http"
)

func main() {
    router := looli.Default()

    router.LoadHTMLGlob("templates/*")
    router.Get("/html", func(c *looli.Context) {
        c.HTML("index.tmpl", looli.JSON{
            "title": "my site",
        })
    })

    http.ListenAndServe(":8080", router)
}
```

templates/index.tmpl

```html
<html>
    <h1>
        {{ .title }}
    </h1>
</html>
```

## Middleware

`looli.Default()` with middleware `Logger()` `Recover()` by default, without middleware use `looli.New()` instead.

### Using middleware

```go
package main

import (
    "net/http"
    "github.com/cssivision/looli"
)

func main() {
    router := looli.New()

    // global middleware
    router.Use(looli.Logger())
    router.Get("/a", func(c *looli.Context) {
        c.Status(200)
        c.String("hello world!\n")
    })

    v1 := router.Prefix("/v1")

    // recover middleware only work for /v1 prefix router
    v1.Use(looli.Recover())
    v1.Get("/a", func(c *looli.Context) {
        panic("error!")
        c.Status(200)
        c.String("hello world!\n")
    })

    http.ListenAndServe(":8080", router)
}
```

### Builtin middlewares

* Logger middleware
* Recover middleware
* [Session middleware](https://github.com/cssivision/looli/tree/master/session)
* [Cors middleware](https://github.com/cssivision/looli/tree/master/cors)
* [Csrf middleware](https://github.com/cssivision/looli/tree/master/csrf)

### Custome middleware

```go
package main

import (
    "log"
    "net/http"
    "github.com/cssivision/looli"
    "time"
)

func Logger() looli.HandlerFunc {
    return func(c *looli.Context) {
        t := time.Now()
        // before request
        c.Next()
        // after request
        latency := time.Since(t)
        log.Print(latency)
    }
}

func main() {
    router := looli.New()

    // global middleware
    router.Use(Logger())
    router.Get("/a", func(c *looli.Context) {
        c.Status(200)
        c.String("hello world!\n")
    })

    http.ListenAndServe(":8080", router)
}
```

# Licenses

All source code is licensed under the [MIT License](https://github.com/cssivision/looli/blob/master/LICENSE).
