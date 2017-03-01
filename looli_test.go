package looli

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewEngine(t *testing.T) {
	engine := New()
	assert.Equal(t, "", engine.basePath)
	assert.Equal(t, engine.engine, engine)
	assert.Equal(t, engine.router, engine.RouterPrefix.router)
	assert.False(t, engine.ForwardedByClientIP)
	assert.Empty(t, engine.Handlers)
}

func TestSetIgnoreCase(t *testing.T) {
	router := New()
	serverResponse := "server response"
	statusCode := 200
	router.SetIgnoreCase(false)
	assert.False(t, router.router.IgnoreCase)
	router.Get("/a/b", func(c *Context) {
		c.Status(statusCode)
		c.String(serverResponse)
	})

	server := httptest.NewServer(router)
	defer server.Close()
	serverURL := server.URL
	resp, err := http.Get(serverURL + "/A/B")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 404, resp.StatusCode)
	resp.Body.Close()

	router.SetIgnoreCase(true)
	assert.True(t, router.router.IgnoreCase)
	resp, err = http.Get(serverURL + "/A/B")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, resp.StatusCode, statusCode)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	assert.Equal(t, string(bodyBytes), serverResponse)
	resp.Body.Close()
}

func TestSetTrailingSlashRedirect(t *testing.T) {
	router := New()
	serverResponse := "server response"
	statusCode := 200
	statusNotFound := 404
	router.Get("/a/b", func(c *Context) {
		c.Status(statusCode)
		c.String(serverResponse)
	})

	server := httptest.NewServer(router)
	defer server.Close()
	serverURL := server.URL
	resp, err := http.Get(serverURL + "/a/b/")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, resp.StatusCode, statusCode)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	assert.Equal(t, string(bodyBytes), serverResponse)
	resp.Body.Close()

	router.SetTrailingSlashRedirect(false)
	resp, err = http.Get(serverURL + "/a/b/")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, resp.StatusCode, statusNotFound)
	resp.Body.Close()
}
