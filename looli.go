package looli

import (
	"github.com/cssivision/router"
	"net/http"
)

type (
	Engine struct {
		// router with basePath
		RouterPrefix

		// router used to match url
		router *router.Router

		// when set true, implements a best effort algorithm to return the real client IP, it parses
		// X-Real-IP and X-Forwarded-For in order to work properly with reverse-proxies such us: nginx or haproxy.
		ForwardedByClientIP bool

		// Enables automatic redirection if the current route can't be matched but a
		// handler for the path with (without) the trailing slash exists.
		// TrailingSlashRedirect: /a/b/ -> /a/b
		// TrailingSlashRedirect: /a/b -> /a/b/
		// default value is true
		TrailingSlashRedirect bool

		// Ignore case when matching URL path.
		IgnoreCase bool
	}
	HandlerFunc func(*Context)
)

func New() *Engine {
	engine := &Engine{
		RouterPrefix:          RouterPrefix{},
		router:                router.New(),
		TrailingSlashRedirect: true,
	}

	engine.RouterPrefix.router = engine.router
	engine.RouterPrefix.router.TrailingSlashRedirect = engine.TrailingSlashRedirect
	engine.RouterPrefix.router.IgnoreCase = false
	engine.RouterPrefix.engine = engine
	return engine
}

// Default return engine instance, add logger, recover handler to it.
func Default() *Engine {
	engine := New()
	engine.RouterPrefix.Handlers = []HandlerFunc{Logger(), Recover()}
	return engine
}

// set IgnoreCase value
func (engine *Engine) SetIgnoreCase(flag bool) {
	engine.RouterPrefix.router.IgnoreCase = flag
}

// set TrailingSlashRedirect value
func (engine *Engine) SetTrailingSlashRedirect(flag bool) {
	engine.RouterPrefix.router.TrailingSlashRedirect = flag
}

// http.Handler interface
func (engine *Engine) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	engine.router.ServeHTTP(rw, req)
}

// short for http.ListenAndServe
func (engine *Engine) Run(address string) (err error) {
	err = http.ListenAndServe(address, engine.router)
	return
}
