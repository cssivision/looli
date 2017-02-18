package looli

import (
	"io"
	"github.com/cssivision/router"
	"math"
	"net/http"
	"net/url"
)

// Context construct Request and ResponseWriter, provide useful methods
type Context struct {
	ResponseWriter
	current  int8

	// http.Request
	Request  *http.Request

	// middleware handler for
	handlers []HandlerFunc

	// Param is a single URL parameter, a map[string]string.
	Params   router.Params

	// Short for Request.URL.String()
	URL      string
}

const abortIndex int8 = math.MaxInt8 / 2

// Next should be used only inside middleware. It executes the pending handlers in the chain
// inside the calling handler
func (c *Context) Next() {
	c.current++
	length := int8(len(c.handlers))
	for ; c.current < length; c.current++ {
		c.handlers[c.current](c)
	}
}

// Abort prevents pending handlers from being called. Note that this will not stop the current handler.
// if you want to stop current handler you should return, after call abort, call Abort to ensure the
// remaining handlers for this request are not called.
func (c *Context) Abort() {
	c.current = abortIndex
}

// IsAborted returns true if the current context was aborted.
func (c *Context) IsAborted() bool {
	return c.current >= abortIndex
}

func (c *Context) Param(name string) string {
	return c.Params[name]
}

// Query returns the keyed url query value if it exists, othewise it returns an empty string `("")`.
// It is shortcut for `c.Request.URL.Query().Get(key)` GET /path?&name=cssivision&age=23
// 		c.Query("name") == "cssivision"
//      c.Query("age") == "23"
// 		c.Query("sex") == ""
func (c *Context) Query(key string) string {
	req := c.Request
	query := req.URL.Query()
	if values, ok := query[key]; ok && len(values) > 0 {
		return values[0]
	}

	return ""
}

// Query returns the keyed url query value if it exists, othewise it returns spectfic defaultValue.
// It is shortcut for `c.Request.URL.Query().Get(key)` GET /path?&name=cssivision&age=23
// 		c.DefaultQuery("name", "balabala") == "cssivision"
//      c.Query("age", "24") == "23"
// 		c.Query("sex", "male") == "male"
func (c *Context) DefaultQuery(key, defaultValue string) string {
	val := c.Query(key)
	if val == "" {
		return defaultValue
	}

	return val
}

// PostForm returns the specified key from a POST urlencoded form or multipart form
// when it exists, otherwise it returns an empty string.
func (c *Context) PostForm(key string) string {
	req := c.Request
	req.ParseForm()
	req.ParseMultipartForm(32 << 20)

	val := ""
	if values := req.PostForm[key]; len(values) > 0 {
		val = values[0]
	}

	if req.MultipartForm != nil && req.MultipartForm.File != nil {
		if values := req.MultipartForm.Value[key]; len(values) > 0 {
			val = values[0]
		}
	}
	return val
}

// DefaultPostForm returns the specified key from a POST urlencoded form or multipart form
// when it exists, otherwise it returns the specified defaultValue string.
func (c *Context) DefaultPostForm(key, defaultValue string) string {
	val := c.PostForm(key)
	if val == "" {
		return defaultValue
	}

	return val
}

func (c *Context) Bind(interface{}) {

}

// Write StatusCode to Response Header
func (c *Context) Status(code int) {
	c.ResponseWriter.WriteHeader(code)
}

// Redirect to location and use http.StatusFound status code
func (c *Context) Redirect(location string) {
	http.Redirect(c.ResponseWriter, c.Request, location, http.StatusFound)
}

func (c *Context) ServeFile(filepath string) {
	http.ServeFile(c.ResponseWriter, c.Request, filepath)
}

// Set Header by key and value
func (c *Context) SetHeader(key, value string) {
	if value == "" {
		c.ResponseWriter.Header().Del(key)
	} else {
		c.ResponseWriter.Header().Set(key, value)
	}
}

// Get Header by key
func (c *Context) Header(key string) string {
	return c.Request.Header.Get(key)
}

// Cookie get cookie from request header by name, if err != nil, return "", err
func (c *Context) Cookie(name string) (string, error) {
	cookie, err := c.Request.Cookie(name)
	if err != nil {
		return "", err
	}

	val, _ := url.QueryUnescape(cookie.Value)
	return val, nil
}

// SetCookie use http.SetCookie to set set-cookie header
func (c *Context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(c.ResponseWriter, cookie)
}

func (c *Context) Pipe(dst io.Writer) {
	io.Copy(dst, c.Request.Body)
}

func (c *Context) ContentType() string {
	if values, _ := c.Request.Header["Content-Type"]; len(values) > 0 {
		return values[0]
	}
	return ""
}

func (c *Context) String() {

}

func (c *Context) JSON() {

}

func (c *Context) HTML() {

}

func (c *Context) Render() {

}
