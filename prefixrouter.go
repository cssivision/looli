package looli

import (
	"github.com/cssivision/router"
	"net/http"
)

type RouterPrefix struct {
	basePath string
	router   *router.Router
}

func (p *RouterPrefix) combineHandlers() {

}

func (p *RouterPrefix) Use(middleware ...HandlerFunc) {
	if len(middleware) == 0 {
		panic("there must be at least one middleware")
	}
}

func (p *RouterPrefix) Get(pattern string, handlers ...HandlerFunc) {
	p.Handle(http.MethodGet, pattern, handlers...)
}

func (p *RouterPrefix) Post(pattern string, handlers ...HandlerFunc) {
	p.Handle(http.MethodGet, pattern, handlers...)
}

func (p *RouterPrefix) Prefix(basePath string) *RouterPrefix {
	return &RouterPrefix{
		basePath: basePath,
		router:   p.router,
	}
}

func (p *RouterPrefix) Handle(method, pattern string, handlers ...HandlerFunc) {
	if len(handlers) == 0 {
		panic("there must be at least one handler")
	}

	handler := compose(handlers...)
	p.router.Handle(method, pattern, handler)
}
