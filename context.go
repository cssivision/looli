package looli

import (
	"net/http"
)

type Context struct {
	current  int
	Request  *http.Request
	Response ResponseWriter
    handlers []HandlerFunc
}

func (c *Context) Next() {
    c.current++
    for ; c.current < len(c.handlers); c.current++ {
        c.handlers[c.current](c)
    }
}
