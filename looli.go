package looli

import (
	"github.com/cssivision/router"
	"net/http"
)

type (
	Engine struct {
		RouterPrefix
		router              *router.Router
		ForwardedByClientIP bool
	}
	HandlerFunc func(*Context)
)

func New() *Engine {
	engine := &Engine{
		RouterPrefix: RouterPrefix{
			basePath: "",
		},
		router: router.New(),
	}

	engine.RouterPrefix.router = engine.router
	engine.RouterPrefix.engine = engine
	return engine
}

func Default() *Engine {
	engine := New()
	engine.RouterPrefix.Handlers = []HandlerFunc{Logger(), Recover()}
	return engine
}

func (engine *Engine) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	engine.router.ServeHTTP(rw, req)
}

func (engine *Engine) Run(address string) (err error) {
	err = http.ListenAndServe(address, engine.router)
	return
}
