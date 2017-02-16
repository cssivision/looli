package looli

import (
	"github.com/cssivision/router"
	"net/http"
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

func (c *Context) Query(key string) string {
	req := c.Request
	query := req.URL.Query()
	if values, ok := query[key]; ok && len(values) > 0 {
		return values[0]
	}

	return ""
}

func (c *Context) DefaultQuery(key, defaultValue string) string {
	val := c.Query(key)
	if val == "" {
		return defaultValue
	}

	return val
}

func (c *Context) PostForm(key string) string {
	req := c.Request
	req.ParseForm()
	req.ParseMultipartForm(32 << 20)
	return ""
}

func (c *Context) DefaultPostForm(key, defaultValue string) string {
	val := c.PostForm(key)
	if val == "" {
		return defaultValue
	}

	return val
}

func (c *Context) Status(code int) {
	c.WriteHeader(code)
}

func (c *Context) Header(key, value string) {
	if value == "" {
		c.ResponseWriter.Header().Del(key)
	} else {
		c.ResponseWriter.Header().Set(key, value)
	}
}

func (c *Context) Redirect(code int, location string) {

}

func (c *Context) ServeFile(filepath string) {
	http.ServeFile(c.ResponseWriter, c.Request, filepath)
}
