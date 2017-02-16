package looli

import (
	"net/http"
	"github.com/cssivision/router"
)

type Context struct {
	ResponseWriter 
	current  int
	Request  *http.Request
    handlers []HandlerFunc
    Params   router.Params
}

func (c *Context) Next() {
    c.current++
    for ; c.current < len(c.handlers); c.current++ {
        c.handlers[c.current](c)
    }
}
