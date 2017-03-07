package looli

import (
	"github.com/cssivision/router"
	"html/template"
	"net/http"
	"path"
	"strings"
)

var (
	defaultStatusCode = http.StatusOK
	default404Body    = "404 page not found"
	default405Body    = "405 method not allowed"
)

// RouterPrefix is used internally to configure router, a RouterPrefix is associated with a basePath
// and an array of handlers (middleware)
type RouterPrefix struct {
	basePath    string
	router      *router.Router
	Middlewares []HandlerFunc
	template    *template.Template
	engine      *Engine
	allNoRoute  []HandlerFunc
	allNoMethod []HandlerFunc
	isPrefix    bool
}

// Use adds middleware to the router.
func (p *RouterPrefix) Use(middleware ...HandlerFunc) {
	if len(middleware) == 0 {
		panic("there must be at least one middleware")
	}

	p.Middlewares = append(p.Middlewares, middleware...)

	p.rebuild404Handlers()
	p.rebuild405Handlers()
}

// Use adds handlers as middleware to the router.
func (p *RouterPrefix) UseHandler(handlers ...Handler) {
	var middlwares []HandlerFunc
	for _, handler := range handlers {
		middlwares = append(middlwares, handler.Handle)
	}
	p.Use(middlwares...)
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

	if p.isPrefix {
		handlers = p.combineHandlers(handlers)
	}
	handler := p.composeHandlers(handlers)
	p.router.Handle(method, pattern, handler)
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

// combine middleware and handlers for specific route
func (p *RouterPrefix) combineHandlers(handlers []HandlerFunc) []HandlerFunc {
	finalSize := len(p.Middlewares) + len(handlers)
	if finalSize >= int(abortIndex) {
		panic("too many handlers")
	}
	mergedHandlers := make([]HandlerFunc, finalSize)
	copyHandlers(mergedHandlers, p.Middlewares)
	copyHandlers(mergedHandlers[len(p.Middlewares):], handlers)
	return mergedHandlers
}

// compose global middleware for all request
func (p *RouterPrefix) composeMiddlewares() func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		context := NewContext(p, rw, req)
		httpHandler := func(c *Context) {
			p.router.ServeHTTP(c.ResponseWriter, c.Request)
		}

		handlers := p.combineHandlers([]HandlerFunc{httpHandler})
		context.handlers = handlers
		context.Next()
	}
}

// Prefix creates a new router prefix. You should add all the routes that have common
// middlwares or the same path prefix. For example, all the routes that use a common
// middlware could be grouped.
func (p *RouterPrefix) Prefix(basePath string) *RouterPrefix {
	return &RouterPrefix{
		basePath:    basePath,
		router:      p.router,
		Middlewares: p.Middlewares,
		engine:      p.engine,
		isPrefix:    true,
	}
}

func copyHandlers(dst, src []HandlerFunc) {
	for index, val := range src {
		dst[index] = val
	}
}

// Construct handler for specific router
func (p *RouterPrefix) composeHandlers(handlers []HandlerFunc) router.Handle {
	return func(rw http.ResponseWriter, req *http.Request, ps router.Params) {
		context := NewContext(p, rw, req)
		context.handlers = handlers
		context.Params = ps

		context.Next()
	}
}

func (p *RouterPrefix) rebuild404Handlers() {
	p.allNoRoute = p.combineHandlers(nil)
}

func (p *RouterPrefix) rebuild405Handlers() {
	p.allNoMethod = p.combineHandlers(nil)
}

// noMethod use as a default handler for router not allowed
func (p *RouterPrefix) noRoute(rw http.ResponseWriter, req *http.Request) {
	context := NewContext(p, rw, req)
	context.handlers = p.allNoRoute

	context.Status(http.StatusNotFound)
	context.String(default404Body)

	context.Next()
}

// noMethod use as a default handler for Method not allowed
func (p *RouterPrefix) noMethod(rw http.ResponseWriter, req *http.Request) {
	context := NewContext(p, rw, req)
	context.handlers = p.allNoMethod

	context.Status(http.StatusMethodNotAllowed)
	context.String(default405Body)
	context.Next()
}

// NoRoute which is called when no matching route is found. If it is not set, noRoute is used.
func (p *RouterPrefix) NoRoute(handlers ...HandlerFunc) {
	if len(handlers) == 0 {
		panic("there must be at least one handler")
	}

	handler := p.composeHandlers(handlers)
	p.router.NoRoute = http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		handler(rw, req, router.Params{})
	})
}

// NoMethod which is called when method is not registered. If it is not set, noMethod is used.
func (p *RouterPrefix) NoMethod(handlers ...HandlerFunc) {
	if len(handlers) == 0 {
		panic("there must be at least one handler")
	}

	handler := p.composeHandlers(handlers)
	p.router.NoMethod = http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		handler(rw, req, router.Params{})
	})
}

func (p *RouterPrefix) handleRequest(rw http.ResponseWriter, req *http.Request) {
	handler := p.composeMiddlewares()
	handler(rw, req)
}
