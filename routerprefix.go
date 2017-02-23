package looli

import (
	"github.com/cssivision/router"
	"html/template"
	"net/http"
	"path"
	"strings"
)

const defaultStatusCode = http.StatusOK

// RouterPrefix is used internally to configure router, a RouterPrefix is associated with a basePath
// and an array of handlers (middleware)
type RouterPrefix struct {
	basePath string
	router   *router.Router
	Handlers []HandlerFunc
	template *template.Template
	engine   *Engine
}

// Use adds middleware to the router.
func (p *RouterPrefix) Use(middleware ...HandlerFunc) {
	if len(middleware) == 0 {
		panic("there must be at least one middleware")
	}
	p.Handlers = append(p.Handlers, middleware...)
}

// Get is a shortcut for router.Handle("GET", path, handle)
func (p *RouterPrefix) Get(pattern string, handlers ...HandlerFunc) {
	p.Handle(http.MethodGet, pattern, handlers...)
}

// Post is a shortcut for router.Handle("Post", path, handle)
func (p *RouterPrefix) Post(pattern string, handlers ...HandlerFunc) {
	p.Handle(http.MethodPost, pattern, handlers...)
}

// Put is a shortcut for router.Handle("Put", path, handle)
func (p *RouterPrefix) Put(pattern string, handlers ...HandlerFunc) {
	p.Handle(http.MethodPut, pattern, handlers...)
}

// Delete is a shortcut for router.Handle("DELETE", path, handle)
func (p *RouterPrefix) Delete(pattern string, handlers ...HandlerFunc) {
	p.Handle(http.MethodDelete, pattern, handlers...)
}

// Head is a shortcut for router.Handle("HEAD", path, handle)
func (p *RouterPrefix) Head(pattern string, handlers ...HandlerFunc) {
	p.Handle(http.MethodHead, pattern, handlers...)
}

// Options is a shortcut for router.Handle("OPTIONS", path, handle)
func (p *RouterPrefix) Options(pattern string, handlers ...HandlerFunc) {
	p.Handle(http.MethodOptions, pattern, handlers...)
}

// Patch is a shortcut for router.Handle("PATCH", path, handle)
func (p *RouterPrefix) Patch(pattern string, handlers ...HandlerFunc) {
	p.Handle(http.MethodPatch, pattern, handlers...)
}

// Any registers a route that matches all the HTTP methods.
// GET, POST, PUT, PATCH, HEAD, OPTIONS, DELETE, CONNECT, TRACE
func (p *RouterPrefix) Any(pattern string, handlers ...HandlerFunc) {
	p.Handle(http.MethodGet, pattern, handlers...)
	p.Handle(http.MethodPost, pattern, handlers...)
	p.Handle(http.MethodPut, pattern, handlers...)
	p.Handle(http.MethodDelete, pattern, handlers...)
	p.Handle(http.MethodHead, pattern, handlers...)
	p.Handle(http.MethodOptions, pattern, handlers...)
	p.Handle(http.MethodPatch, pattern, handlers...)
	p.Handle(http.MethodTrace, pattern, handlers...)
	p.Handle(http.MethodConnect, pattern, handlers...)
}

// Handle registers a new request handle and middleware with the given path and method.
func (p *RouterPrefix) Handle(method, pattern string, handlers ...HandlerFunc) {
	if len(handlers) == 0 {
		panic("there must be at least one handler")
	}

	if p.basePath != "" {
		pattern = p.basePath + pattern
	}

	handlers = p.combineHandlers(handlers)
	muxHandler := p.composeMiddleware(handlers)
	p.router.Handle(method, pattern, muxHandler)
}

func (p *RouterPrefix) LoadHTMLGlob(pattern string) {
	templ := template.Must(template.ParseGlob(pattern))
	p.SetHTMLTemplate(templ)
}

func (p *RouterPrefix) LoadHTMLFiles(files ...string) {
	templ := template.Must(template.ParseFiles(files...))
	p.SetHTMLTemplate(templ)
}

func (p *RouterPrefix) SetHTMLTemplate(templ *template.Template) {
	p.template = templ
}

// StaticFile register router pattern and response file in path
func (p *RouterPrefix) StaticFile(pattern, filepath string) {
	if strings.Contains(pattern, ":") || strings.Contains(pattern, "*") {
		panic("URL parameters can not be used when serving a static folder")
	}
	handler := func(c *Context) {
		c.ServeFile(filepath)
	}

	p.Head(pattern, handler)
	p.Get(pattern, handler)
}

// Static register router pattern and response file in the request url
func (p *RouterPrefix) Static(pattern, dir string) {
	if strings.Contains(pattern, ":") || strings.Contains(pattern, "*") {
		panic("URL parameters can not be used when serving a static folder")
	}

	fileServer := http.StripPrefix(pattern, http.FileServer(http.Dir(dir)))
	handler := func(c *Context) {
		fileServer.ServeHTTP(c.ResponseWriter, c.Request)
	}

	urlPattern := path.Join(pattern, "/*filepath")
	p.Head(urlPattern, handler)
	p.Get(urlPattern, handler)
}

// Configurable NotFound which is called when no matching route is
// found. If it is not set, http.NotFound is used.
func (p *RouterPrefix) NotFound(handlers ...HandlerFunc) {
	if len(handlers) == 0 {
		panic("there must be at least one handler")
	}

	handler := p.composeMiddleware(handlers)
	p.router.NotFound = http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusNotFound)
		handler(rw, req, router.Params{})
	})
}

// combine middleware and handlers for specific route
func (p *RouterPrefix) combineHandlers(handlers []HandlerFunc) []HandlerFunc {
	finalSize := len(p.Handlers) + len(handlers)
	if finalSize >= int(abortIndex) {
		panic("too many handlers")
	}
	mergedHandlers := make([]HandlerFunc, finalSize)
	copyHandlers(mergedHandlers, p.Handlers)
	copyHandlers(mergedHandlers[len(p.Handlers):], handlers)
	return mergedHandlers
}

// Prefix creates a new router prefix. You should add all the routes that have common
// middlwares or the same path prefix. For example, all the routes that use a common
// middlware could be grouped.
func (p *RouterPrefix) Prefix(basePath string) *RouterPrefix {
	return &RouterPrefix{
		basePath: basePath,
		router:   p.router,
		Handlers: p.Handlers,
	}
}

func copyHandlers(dst, src []HandlerFunc) {
	for index, val := range src {
		dst[index] = val
	}
}

// Construct handler for specific router
func (p *RouterPrefix) composeMiddleware(handlers []HandlerFunc) router.Handle {
	return func(rw http.ResponseWriter, req *http.Request, ps router.Params) {
		context := &Context{
			ResponseWriter: rw,
			Request:        req,
			handlers:       handlers,
			current:        -1,
			Params:         ps,
			Path:           req.URL.String(),
			template:       p.template,
			engine:         p.engine,
			statusCode:     defaultStatusCode,
		}

		context.Next()
	}
}
