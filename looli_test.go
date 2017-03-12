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
	assert.Empty(t, engine.Middlewares)
}

func TestDefault(t *testing.T) {
	router := Default()
	assert.Equal(t, 2, len(router.Middlewares))
}

func TestDeafultnoMethod(t *testing.T) {
	router := New()
	router.Get("/a", func(c *Context) {
		noMethod(c)
	})

	server := httptest.NewServer(router)
	defer server.Close()

	serverURL := server.URL
	resp, err := http.Get(serverURL + "/a")
	assert.Nil(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, default405Body, string(bodyBytes))
}

func TestDeafultnoRoute(t *testing.T) {
	router := New()
	assert.False(t, router.router.IgnoreCase)
	router.Get("/a", func(c *Context) {
		noRoute(c)
	})

	server := httptest.NewServer(router)
	defer server.Close()

	serverURL := server.URL
	resp, err := http.Get(serverURL + "/a")
	assert.Nil(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, default404Body, string(bodyBytes))
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

func TestNoMethod(t *testing.T) {
	t.Run("no method", func(t *testing.T) {
		router := New()
		serverResponse := "server response"
		statusCode := 404

		router.Use(func(c *Context) {
			c.SetHeader("fake-header", "fake")
		})

		router.NoMethod(func(c *Context) {
			c.Status(statusCode)
			c.String(serverResponse)
		})

		server := httptest.NewServer(router)
		defer server.Close()

		serverURL := server.URL
		resp, err := http.Get(serverURL + "/a")
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		assert.Equal(t, statusCode, resp.StatusCode)
		assert.Equal(t, "fake", resp.Header.Get("fake-header"))

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, serverResponse, string(bodyBytes))
	})

	t.Run("no method handler should panic", func(t *testing.T) {
		router := New()
		assert.Panics(t, func() {
			router.NoMethod()
		})
	})
}

func TestNoRoute(t *testing.T) {
	t.Run("no route", func(t *testing.T) {
		router := New()
		serverResponse := "server response"
		statusCode := 404
		router.Use(func(c *Context) {
			c.SetHeader("fake-header", "fake")
		})

		router.NoRoute(func(c *Context) {
			c.Status(statusCode)
			c.String(serverResponse)
		})
		router.Get("/a/b", func(c *Context) {})
		server := httptest.NewServer(router)
		defer server.Close()

		serverURL := server.URL
		resp, err := http.Get(serverURL + "/a")
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		assert.Equal(t, statusCode, resp.StatusCode)
		assert.Equal(t, "fake", resp.Header.Get("fake-header"))

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, serverResponse, string(bodyBytes))
	})

	t.Run("no route handler should panic", func(t *testing.T) {
		router := New()
		assert.Panics(t, func() {
			router.NoRoute()
		})
	})
}
