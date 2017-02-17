package looli

import (
	"github.com/cssivision/router"
	"net/http"
)

// RouterPrefix is used internally to configure router, a RouterPrefix is associated with a basePath
// and an array of handlers (middleware)
type RouterPrefix struct {
	basePath string
	router   *router.Router
	Handlers []HandlerFunc
}

// Use adds middleware to the router.
func (p *RouterPrefix) Use(middleware ...HandlerFunc) {
	if len(middleware) == 0 {
		panic("there must be at least one middleware")
	}
	p.Handlers = middleware
}

// Get is a shortcut for router.Handle("GET", path, handle)
func (p *RouterPrefix) Get(pattern string, handlers ...HandlerFunc) {
	p.Handle(http.MethodGet, pattern, handlers...)
}

// Post is a shortcut for router.Handle("Post", path, handle)
func (p *RouterPrefix) Post(pattern string, handlers ...HandlerFunc) {
	p.Handle(http.MethodGet, pattern, handlers...)
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

// Trace is a shortcut for router.Handle("TRACE", path, handle)
func (p *RouterPrefix) Trace(pattern string, handlers ...HandlerFunc) {
	p.Handle(http.MethodTrace, pattern, handlers...)
}

// Options is a shortcut for router.Handle("OPTIONS", path, handle)
func (p *RouterPrefix) Options(pattern string, handlers ...HandlerFunc) {
	p.Handle(http.MethodOptions, pattern, handlers...)
}

// Patch is a shortcut for router.Handle("PATCH", path, handle)
func (p *RouterPrefix) Patch(pattern string, handlers ...HandlerFunc) {
	p.Handle(http.MethodPatch, pattern, handlers...)
}

func (p *RouterPrefix) combineHandlers(handlers []HandlerFunc) []HandlerFunc {
	finalSize := len(p.Handlers) + len(handlers)
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

// Handle registers a new request handle and middleware with the given path and method.
func (p *RouterPrefix) Handle(method, pattern string, handlers ...HandlerFunc) {
	if len(handlers) == 0 {
		panic("there must be at least one handler")
	}

	if p.basePath != "" {
		pattern = p.basePath + pattern
	}

	handlers = p.combineHandlers(handlers)
	muxHandler := composeMiddleware(handlers)
	p.router.Handle(method, pattern, muxHandler)
}

func copyHandlers(dst, src []HandlerFunc) {
	for index, val := range src {
		dst[index] = val
	}
}

// Construct handler for specific router
func composeMiddleware(handlers []HandlerFunc) router.Handle {
	return func(rw http.ResponseWriter, req *http.Request, ps router.Params) {
		context := &Context{
			Request:  req,
			handlers: handlers,
			current:  -1,
			ResponseWriter: ResponseWriter{
				ResponseWriter: rw,
			},
			Params: ps,
		}

		context.Next()
	}
}
