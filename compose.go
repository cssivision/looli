package looli

import (
    "net/http"
    "github.com/cssivision/router"
)

func compose(handlers ...HandlerFunc) router.Handle {
    return func(rw http.ResponseWriter, req *http.Request, ps router.Params) {
    }
}