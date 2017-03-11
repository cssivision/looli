package looli

import (
	"net/http"
)

type (
	Engine struct {
		// router with basePath
		RouterPrefix

		// router used to match url
		router *Router

		// when set true, implements a best effort algorithm to return the real client IP, it parses
		// X-Real-IP and X-Forwarded-For in order to work properly with reverse-proxies such us: nginx or haproxy.
		ForwardedByClientIP bool
	}
	HandlerFunc func(*Context)
)

type Handler interface {
	Handle(*Context)
}

func New() *Engine {
	engine := &Engine{
		RouterPrefix: RouterPrefix{},
		router:       NewRouter(),
	}

	engine.RouterPrefix.engine = engine
	engine.RouterPrefix.router = engine.router
	engine.router.IgnoreCase = false
	engine.router.TrailingSlashRedirect = true
	engine.router.NoRoute = []HandlerFunc{noRoute}
	engine.router.NoMethod = []HandlerFunc{noMethod}
	return engine
}

// noRoute use as a default handler for router not matched
func noRoute(c *Context) {
	c.Status(http.StatusNotFound)
	c.String(default404Body)
}

// noMethod use as a default handler for Method not allowed
func noMethod(c *Context) {
	c.Status(http.StatusMethodNotAllowed)
	c.String(default405Body)
}

// Default return engine instance, add logger, recover handler to it.
func Default() *Engine {
	engine := New()
	engine.RouterPrefix.Use(Logger(), Recover())
	return engine
}

// NoRoute which is called when no matching route is found. If it is not set, noRoute is used.
func (engine *Engine) NoRoute(handlers ...HandlerFunc) {
	if len(handlers) == 0 {
		panic("there must be at least one handler")
	}

	engine.router.NoRoute = handlers
}

// NoMethod which is called when method is not registered. If it is not set, noMethod is used.
func (engine *Engine) NoMethod(handlers ...HandlerFunc) {
	if len(handlers) == 0 {
		panic("there must be at least one handler")
	}

	engine.router.NoMethod = handlers
}

// set IgnoreCase value
func (engine *Engine) SetIgnoreCase(ignoreCase bool) {
	engine.router.IgnoreCase = ignoreCase
}

// set TrailingSlashRedirect value
func (engine *Engine) SetTrailingSlashRedirect(redirect bool) {
	engine.router.TrailingSlashRedirect = redirect
}

// http.Handler interface
func (engine *Engine) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// engine.router.ServeHTTP(rw, req)
	engine.RouterPrefix.handleRequest(rw, req)
}
