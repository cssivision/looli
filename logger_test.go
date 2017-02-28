package looli

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogger(t *testing.T) {
	buffer := new(bytes.Buffer)
	router := New()
	router.Use(LoggerWithWriter(buffer))
	router.Get("/a", func(c *Context) {})
	router.Post("/a", func(c *Context) {})
	router.Put("/a", func(c *Context) {})
	router.Delete("/a", func(c *Context) {})
	router.Patch("/a", func(c *Context) {})
	router.Head("/a", func(c *Context) {})
	router.Options("/a", func(c *Context) {})

	methods := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch, http.MethodHead, http.MethodOptions}
	for _, method := range methods {
		buffer.Reset()
		issueRequest(t, router, method, "/a")
		assert.Contains(t, buffer.String(), method)
		assert.Contains(t, buffer.String(), "/a")
		assert.Contains(t, buffer.String(), "200")

		buffer.Reset()
		issueRequest(t, router, method, "/a/b")
		assert.Contains(t, buffer.String(), method)
		assert.Contains(t, buffer.String(), "/a/b")
		assert.Contains(t, buffer.String(), "404")
	}
}

func issueRequest(t *testing.T, router *Engine, method, path string) {
	server := httptest.NewServer(router)
	defer server.Close()

	serverURL := server.URL
	getReq, err := http.NewRequest(method, serverURL+path, nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(getReq)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
}
