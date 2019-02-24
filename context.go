package looli

import (
	"html/template"
	"math"
	"net"
	"net/http"
	"net/url"
	"strings"
)

var defaultStatusCode = http.StatusOK

// Context construct Request and ResponseWriter, provide useful methods
type Context struct {
	http.ResponseWriter

	// current handler that processing request
	current int8

	// Short for http.Request
	Request *http.Request

	// middleware handlers
	handlers []HandlerFunc

	// Param is URL parameter, a map[string]string.
	Params Params

	// Short for Request.URL.Path
	Path string

	// Short for Request.Method
	Method string

	// templete is use to render HTML
	template *template.Template
	engine   *Engine

	// statusCode that write to response
	statusCode int

	// Error when processing request
	Err *Error
}

type JSON map[string]interface{}

const abortIndex int8 = math.MaxInt8 / 2

func NewContext(p *RouterPrefix, rw http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		ResponseWriter: rw,
		Request:        req,
		current:        -1,
		Path:           req.URL.Path,
		Method:         req.Method,
		template:       p.engine.Template,
		engine:         p.engine,
		statusCode:     defaultStatusCode,
	}
}

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

// AbortWithStatus prevents pending handlers from being called and set statuscode. Note that this will not
// stop the current handler. if you want to stop current handler you should return, after call abort, call
// Abort to ensure the remaining handlers for this request are not called.
func (c *Context) AbortWithStatus(code int) {
	c.statusCode = code
	c.Status(code)
	c.Abort()
}

// IsAborted returns true if the current context was aborted.
func (c *Context) IsAborted() bool {
	return c.current >= abortIndex
}

// Param return the parameters by name in the request path
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

// ClientIP implements a best effort algorithm to return the real client IP, it parses
// X-Real-IP and X-Forwarded-For in order to work properly with reverse-proxies such us: nginx or haproxy.
func (c *Context) ClientIP() string {
	if c.engine.ForwardedByClientIP {
		clientIP := strings.TrimSpace(c.Header("X-Real-Ip"))
		if len(clientIP) > 0 {
			return clientIP
		}
		clientIP = c.Header("X-Forwarded-For")
		if index := strings.IndexByte(clientIP, ','); index >= 0 {
			clientIP = clientIP[0:index]
		}
		clientIP = strings.TrimSpace(clientIP)
		if len(clientIP) > 0 {
			return clientIP
		}
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(c.Request.RemoteAddr)); err == nil {
		return ip
	}
	return ""
}

// Bind checks the Content-Type to select a binding engine automatically,
// Depending the "Content-Type" header different bindings are used:
// 		"application/json" --> JSON
// 		"application/xml"  --> XML
// otherwise --> returns an error
// It parses the request's body as JSON if Content-Type == "application/json" using JSON or XML as a JSON input.
// It decodes the json payload into the struct specified as a pointer.
// Like ParseBody() but this method also writes a 400 error if the json is not valid.
func (c *Context) Bind(data BindingStruct) error {
	binding := bindDefault(c.Request.Method, c.ContentType())
	if err := binding.Bind(c.Request, data); err != nil {
		return err
	}
	return data.Validate()
}

// WriteHeader sends an HTTP response header with status code.
// If WriteHeader is not called explicitly, the first call to Write
// will trigger an implicit WriteHeader(http.StatusOK).
// Thus explicit calls to WriteHeader are mainly used to
// send error codes.
func (c *Context) Status(code int) {
	c.statusCode = code
	c.ResponseWriter.WriteHeader(code)
}

// Redirect replies to the request with a redirect to url, which may be a path relative to the request path.
func (c *Context) Redirect(location string) {
	http.Redirect(c.ResponseWriter, c.Request, location, http.StatusFound)
}

// ServeFile replies to the request with the contents of the named file or directory.
// If the provided file or directory name is a relative path, it is interpreted
// relative to the current directory and may ascend to parent directories. If
// the provided name is constructed from user input, it should be sanitized
// before calling ServeFile. As a precaution, ServeFile will reject requests
// where r.URL.Path contains a ".." path element.

// As a special case, ServeFile redirects any request where r.URL.Path ends in
// "/index.html" to the same path, without the final "index.html". To avoid
// such redirects either modify the path or use ServeContent.
func (c *Context) ServeFile(filepath string) {
	http.ServeFile(c.ResponseWriter, c.Request, filepath)
}

// Get gets the first value associated with the given key. It is case insensitive
func (c *Context) Header(key string) string {
	return c.Request.Header.Get(key)
}

// Set sets the header entries associated with key to the single element value.
// It replaces any existing values associated with key.
func (c *Context) SetHeader(key, value string) {
	if value == "" {
		c.ResponseWriter.Header().Del(key)
	} else {
		c.ResponseWriter.Header().Set(key, value)
	}
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

// ContentType return content-type from header
func (c *Context) ContentType() string {
	if values, _ := c.Request.Header["Content-Type"]; len(values) > 0 {
		return values[0]
	}
	return ""
}

func (c *Context) Error(err error) {
	var parsedError *Error
	switch err.(type) {
	case *Error:
		parsedError = err.(*Error)
	default:
		parsedError = &Error{
			Err: err,
		}
	}

	c.Err = parsedError
}

// String write format string to response
func (c *Context) String(format string, values ...interface{}) {
	if err := renderString(c.ResponseWriter, format, values...); err != nil {
		panic(err)
	}
}

// JSON write obj to response
func (c *Context) JSON(data interface{}) {
	if err := renderJSON(c.ResponseWriter, data); err != nil {
		panic(err)
	}
}

// SetResult set response code and msg
func (c *Context) SetResult(code int, msg string) {
	data := map[string]interface{}{
		"code": code,
		"msg":  msg,
	}

	if err := renderJSON(c.ResponseWriter, data); err != nil {
		panic(err)
	}
}

// SetBody return json body
func (c *Context) SetBody(data interface{}) {
	rsp := map[string]interface{}{
		"code": 0,
		"msg":  "ok",
		"data": data,
	}

	if err := renderJSON(c.ResponseWriter, rsp); err != nil {
		panic(err)
	}
}

// HTML rendder html resp
func (c *Context) HTML(name string, data interface{}) {
	if err := renderHTML(c.ResponseWriter, c.template, name, data); err != nil {
		panic(err)
	}
}
