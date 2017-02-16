package looli

import (
	"github.com/cssivision/router"
	"net/http"
)

type RouterPrefix struct {
	basePath string
	router   *router.Router
	Handlers []HandlerFunc
}

func (p *RouterPrefix) Use(middleware ...HandlerFunc) {
	if len(middleware) == 0 {
		panic("there must be at least one middleware")
	}
	p.Handlers = middleware
}

func (p *RouterPrefix) Get(pattern string, handlers ...HandlerFunc) {
	p.Handle(http.MethodGet, pattern, handlers...)
}

func (p *RouterPrefix) Post(pattern string, handlers ...HandlerFunc) {
	p.Handle(http.MethodGet, pattern, handlers...)
}

func (p *RouterPrefix) Put(pattern string, handlers ...HandlerFunc) {
	p.Handle(http.MethodPut, pattern, handlers...)
}

func (p *RouterPrefix) Delete(pattern string, handlers ...HandlerFunc) {
	p.Handle(http.MethodDelete, pattern, handlers...)
}

func (p *RouterPrefix) Head(pattern string, handlers ...HandlerFunc) {
	p.Handle(http.MethodHead, pattern, handlers...)
}

func (p *RouterPrefix) Trace(pattern string, handlers ...HandlerFunc) {
	p.Handle(http.MethodTrace, pattern, handlers...)
}

func (p *RouterPrefix) Options(pattern string, handlers ...HandlerFunc) {
	p.Handle(http.MethodOptions, pattern, handlers...)
}

func (p *RouterPrefix) combineHandlers(handlers []HandlerFunc) []HandlerFunc {
	finalSize := len(p.Handlers) + len(handlers)
	mergedHandlers := make([]HandlerFunc, finalSize)
	copyHandlers(mergedHandlers, p.Handlers)
	copyHandlers(mergedHandlers[len(p.Handlers):], handlers)
	return mergedHandlers
}

func (p *RouterPrefix) Prefix(basePath string) *RouterPrefix {
	return &RouterPrefix{
		basePath: basePath,
		router:   p.router,
		Handlers: p.Handlers,
	}
}

func (p *RouterPrefix) Handle(method, pattern string, handlers ...HandlerFunc) {
	if len(handlers) == 0 {
		panic("there must be at least one handler")
	}

	if p.basePath != "" {
		pattern = p.basePath + pattern
	}

	handlers = p.combineHandlers(handlers)
	muxHandler := compose(handlers)
	p.router.Handle(method, pattern, muxHandler)
}

func copyHandlers(dst, src []HandlerFunc) {
	for index, val := range src {
		dst[index] = val
	}
}

func compose(handlers []HandlerFunc) router.Handle {
    return func(rw http.ResponseWriter, req *http.Request, ps router.Params) {
    	context := &Context{
    		Request: req,
    		handlers: handlers,
    		current: -1,
    		ResponseWriter: ResponseWriter{
    			ResponseWriter: rw,
    		},
    		Params: ps,
    	}

    	context.Next()
    }
}
