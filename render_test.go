package looli

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetContentType(t *testing.T) {
	router := New()
	router.Get("/", func(c *Context) {
		setContentType(c.ResponseWriter, plainContentType)
	})

	server := httptest.NewServer(router)
	defer server.Close()

	serverURL := server.URL
	resp, err := http.Get(serverURL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, plainContentType[0], resp.Header.Get("Content-Type"))
}

func TestRenderString(t *testing.T) {
	router := New()
	router.Get("/", func(c *Context) {
		renderString(c.ResponseWriter, "hello ")
		renderString(c.ResponseWriter, "world %v!", "cssivision")
	})

	server := httptest.NewServer(router)
	defer server.Close()

	serverURL := server.URL
	resp, err := http.Get(serverURL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, plainContentType[0], resp.Header.Get("Content-Type"))
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "hello world cssivision!", string(bodyBytes))
}

func TestRenderJSON(t *testing.T) {
	type Info struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	router := New()
	router.Get("/", func(c *Context) {
		renderJSON(c.ResponseWriter, JSON{
			"name": "cssivision",
			"age":  21,
		})
	})

	server := httptest.NewServer(router)
	defer server.Close()

	serverURL := server.URL
	resp, err := http.Get(serverURL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, jsonContentType[0], resp.Header.Get("Content-Type"))
	info := new(Info)
	json.NewDecoder(resp.Body).Decode(info)
	assert.Equal(t, "cssivision", info.Name)
	assert.Equal(t, 21, info.Age)
}

func TestRenderHTML(t *testing.T) {
	t.Run("render with name", func(t *testing.T) {
		router := New()
		router.LoadHTMLGlob("test/templates/*")
		router.Get("/", func(c *Context) {
			err := renderHTML(c.ResponseWriter, c.template, "index.tmpl", JSON{
				"title": "Posts",
			})
			assert.Nil(t, err)
		})

		server := httptest.NewServer(router)
		defer server.Close()

		serverURL := server.URL
		resp, err := http.Get(serverURL)
		assert.Nil(t, err)
		defer resp.Body.Close()

		assert.Equal(t, htmlContentType[0], resp.Header.Get("Content-Type"))
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		assert.Nil(t, err)
		assert.True(t, strings.Contains(string(bodyBytes), "Posts"))
	})

	t.Run("render without name", func(t *testing.T) {
		router := New()
		router.LoadHTMLGlob("test/templates/*")
		router.Get("/", func(c *Context) {
			err := renderHTML(c.ResponseWriter, c.template, "", JSON{
				"title": "Posts",
			})
			assert.Nil(t, err)
		})

		server := httptest.NewServer(router)
		defer server.Close()

		serverURL := server.URL
		resp, err := http.Get(serverURL)
		assert.Nil(t, err)
		defer resp.Body.Close()

		assert.Equal(t, htmlContentType[0], resp.Header.Get("Content-Type"))
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		assert.Nil(t, err)

		assert.True(t, strings.Contains(string(bodyBytes), "Posts"))
	})
}
