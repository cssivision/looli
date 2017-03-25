# session

session provides cookie session operation.

The key features are: 
* Simple API: use it as an easy way to set signed (and optionally encrypted) cookies.
* Built-in backends to store sessions in cookies
* Authentication and encryption.
* Multiple sessions per request

## Usage

store value in cookie session.

```go
package main

import (
	"github.com/cssivision/looli"
	"github.com/cssivision/looli/session"
	"net/http"
	"fmt"
)

func main() {
	router := looli.Default()

	// secret used to signed cookies, 
	secret := "secret"

	// aesKey used to encrypted the session values, the AES key, either 16, 24, 
	// or 32 bytes to select AES-128, AES-192, or AES-256, if is empty, will 
	// not encrypted session values.
	aesKey := "1111111111111111"
	sessions := session.NewSessions(secret, aesKey)

	router.Get("/a", func(ctx *looli.Context) {
		// Get a session. ignoring the error resulted from decoding an existing 
		// session: Get() always returns a session, even if empty.
		session, err := sessions.Get(ctx, "sess")
		if err != nil {
			fmt.Println("err: ", err)
			return
		}

		// Values is type map[interface{}]interface{}
		if session.Values["looli"] == nil {
			session.Values["looli"] = 1
		} else {
			session.Values["looli"] = session.Values["looli"].(int) + 1
		}

		// Save it before we write to the response/return from the handler.
		err = session.Save(ctx)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(session.Values["looli"])
		ctx.String("Hello World!")
	})

	http.ListenAndServe(":8080", router)
}
```

## Licenses

All source code is licensed under the [MIT License](https://github.com/cssivision/looli/blob/master/LICENSE).
