package looli

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestNewRouter(t *testing.T) {
	router := NewRouter()
	assert.True(t, router.TrailingSlashRedirect)
	assert.NotNil(t, router.allowMethods)
	assert.NotNil(t, router.tree)
	assert.NotNil(t, router.tree.children)
	assert.NotNil(t, router.tree.handlers)
}

func TestAddhandlers(t *testing.T) {
	router := NewRouter()
	router.Handle(http.MethodGet, "/", []HandlerFunc{func(c *Context) {}})
	assert.Panics(t, func() {
		router.Handle(http.MethodGet, "/", []HandlerFunc{func(c *Context) {}})
	})
}
