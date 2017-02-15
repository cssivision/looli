package looli

import (
    "net/http"
	"github.com/cssivision/router"
)

type HandlerFunc func(*Context)

type Engine struct {
	*router.Router
    *PrefixRouter
}

func New() *Engine {
	engine := &Engine{
		Router: router.New(),
	}

    return engine
}

func (engine *Engine) Use(middleware ...HandlerFunc) {

}

func (engine *Engine) Get(pattern string, handlers ...HandlerFunc) {
    engine.Handle(http.MethodGet, pattern, handlers...)
}

func (engine *Engine) Post(pattern string, handlers ...HandlerFunc) {
    engine.Handle(http.MethodGet, pattern, handlers...)
}

func (engine *Engine) Prefix(basePath string) {

}

func (engine *Engine) Handle(method, pattern string, handlers ...HandlerFunc) {
    if len(handlers) == 0 {
        panic("there must be at least one handler")
    }

    handler := compose(handlers...)
    engine.Router.Handle(method, pattern, handler)
}

func (engine *Engine) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
    engine.Router.ServeHTTP(rw, req)
}
