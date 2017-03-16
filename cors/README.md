# cors

a cors middleware for looli.

## Usage 

```go
package main

import (
	"github.com/cssivision/looli"
	"github.com/cssivision/looli/cors"
	"net/http"
)

func main() {
	router := looli.Default()

	// Default return cors middleware with default options being all origins accepted, 
	// AllowMethods default is []string{"GET", "HEAD", "PUT", "POST", "DELETE", "PATCH"}.
	router.Use(cors.Default())

	router.Get("/a", func(ctx *looli.Context) {
		ctx.String("cors response!\n")
	})

	http.ListenAndServe(":8080", router)
}
```

The server now runs on localhost:8080:

```
curl -D - -H 'Origin: http://looli.xyz' http://localhost:8080/a

HTTP/1.1 200 OK
Access-Control-Allow-Origin: http://looli.xyz
Content-Type: text/plain; charset=utf-8
Vary: Origin
Date: Thu, 16 Mar 2017 03:34:11 GMT
Content-Length: 15

cors response!
```

## Parameters

Parameters are passed to the middleware the cors.New method as follow:

```go
corsMiddleware := cors.New(cors.Options{
    AllowedOrigins: []string{"*"},
    AllowCredentials: true,
})
```

* **AllowOrigins** `[]string`: AllowedOrigins is a list of origins a cross-domain request can be executed from, If the special "*" value is present in the list, all origins will be allowed.

* **AllowOriginsFunc** `func(string) bool`: AllowOriginFunc is a custom function to validate the origin. It take the origin as argument and returns true if allowed or false otherwise. If this option is set, the content of AllowedOrigins is ignored.

* **AllowMethods** `[]string`: specifies the method or methods allowed when accessing the resource. This is used in response to a preflight request, default is []string{"GET", "HEAD", "PUT", "POST", "DELETE", "PATCH"}.

* **AllowHeaders** `[]string`: header is used in response to a preflight request to indicate which HTTP headers can be used when making the actual request.

* **AllowCredentials** `bool`: The Access-Control-Allow-Credentials header Indicates whether or not the response to the request can be exposed when the credentials flag is true.  When used as part of a response to a preflight request, this indicates whether or not the actual request can be made using credentials. Note that simple GET requests are not preflighted, and so if a request is made for a resource with credentials, if this header is not returned with the resource, the response is ignored by the browser and not returned to web content.

* **ExposeHeaders** `[]string`: ExposeHeaders header lets a server whitelist headers that browsers are allowed to access. For example: Access-Control-Expose-Headers: "X-My-Custom-Header, X-Another-Custom-Header" This allows the X-My-Custom-Header and X-Another-Custom-Header headers to be exposed to the browser.

* **MaxAge** `time.Duration`: MaxAge indicates how long the results of a preflight request can be cached.

## Licenses

All source code is licensed under the [MIT License](https://github.com/cssivision/looli/blob/master/LICENSE).