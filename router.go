package looli

import (
	"net/http"
	"strings"
)

// Router is a http.Handler which can be used to dispatch requests to different
// handler functions via configurable routes
type Router struct {
	// tree used to keep handler with path
	tree *node

	// Ignore case when matching URL path.
	IgnoreCase bool

	// Enables automatic redirection if the current route can't be matched but a
	// handler for the path with (without) the trailing slash exists.
	// TrailingSlashRedirect: /a/b/ -> /a/b
	// TrailingSlashRedirect: /a/b -> /a/b/
	TrailingSlashRedirect bool

	// Configurable http.Handler which is called when no matching route is
	// found. If it is not set, http.NotFound is used.
	NoRoute []HandlerFunc

	// Configurable http.Handler which is called when method is not allowed. If it is not set, http.NotFound is used.
	NoMethod []HandlerFunc

	// Methods which has been registered
	allowMethods map[string]bool
}

// Handle is a function that can be registered to a route to handle HTTP
// requests. Like http.HandlerFunc, but has a third parameter for the
// values of named/wildcards parameters.
type Handle func(http.ResponseWriter, *http.Request, Params)

// Param is a single URL parameter, a map[string]string.
type Params map[string]string

// New returns a new initialized Router, with default configuration
func NewRouter() *Router {
	router := &Router{
		tree: &node{
			children: make(map[string]*node),
			handlers: make(map[string][]HandlerFunc),
		},
		TrailingSlashRedirect: true,
		allowMethods:          make(map[string]bool),
	}
	return router
}

// Handle registers a new request handle with the given path and method.
// For GET, POST, PUT, PATCH and DELETE requests the respective shortcut
// functions can be used.
func (r *Router) Handle(method, pattern string, handlers []HandlerFunc) {
	if pattern[0] != '/' {
		panic("path must begin with '/', '" + pattern + "'")
	}

	if method == "" {
		panic("invalid http method")
	}

	if r.tree == nil {
		r.tree = &node{
			children: make(map[string]*node),
			handlers: make(map[string][]HandlerFunc),
		}
	}

	if r.IgnoreCase {
		pattern = strings.ToLower(pattern)
	}

	if !r.allowMethods[method] {
		r.allowMethods[method] = true
	}
	r.tree.insert(pattern).addHandlers(method, handlers)
}

func (r *Router) handleRequest(c *Context) {
	rw := c.ResponseWriter
	req := c.Request

	pattern := req.URL.Path
	if r.IgnoreCase {
		pattern = strings.ToLower(pattern)
	}

	// handle for matched request
	n, ps, tsr := r.tree.find(pattern)
	if n != nil {
		if handlers := n.handlers[req.Method]; handlers != nil {
			c.handlers = append(c.handlers, handlers...)
			c.Params = ps
			c.Next()
			return
		}
	} else {
		// handle for trailing slash redirect
		if r.TrailingSlashRedirect && tsr {
			path := req.URL.Path
			if len(path) > 1 && path[len(path)-1] == '/' {
				pattern = path[:len(path)-1]
			} else {
				pattern = path + "/"
			}

			http.Redirect(rw, req, pattern, http.StatusMovedPermanently)
			return
		}
	}

	if !r.allowMethods[req.Method] {
		if r.NoMethod != nil {
			c.handlers = append(c.handlers, r.NoMethod...)
			c.Params = ps
			c.Next()
		} else {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			rw.Write([]byte(default405Body))
		}
		return
	}

	if r.NoRoute != nil {
		c.handlers = append(c.handlers, r.NoRoute...)
		c.Params = ps
		c.Next()
	} else {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte(default404Body))
	}
}
