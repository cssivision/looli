package looli

import (
	"github.com/cssivision/router"
	"net/http"
)

type HandlerFunc func(*Context)

type Engine struct {
	RouterPrefix
	router *router.Router
}

func New() *Engine {
	engine := &Engine{
		router: router.New(),
		RouterPrefix: RouterPrefix{
			basePath: "",
		},
	}

	engine.RouterPrefix.router = engine.router
	return engine
}

func (engine *Engine) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	engine.router.ServeHTTP(rw, req)
}
