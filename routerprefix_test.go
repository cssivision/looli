package looli

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMethod(t *testing.T) {
	handleMethod(http.MethodGet, t)
	handleMethod(http.MethodOptions, t)
	handleMethod(http.MethodPatch, t)
	handleMethod(http.MethodDelete, t)
	handleMethod(http.MethodTrace, t)
	handlePostPutMethod(http.MethodPost, t)
	handlePostPutMethod(http.MethodPut, t)
}

func TestHeadMethod(t *testing.T) {
	serverResponse := "server response"
	statusCode := 404
	router := New()
	router.Head("/a/b", func(c *Context) {
		c.Status(statusCode)
	})

	server := httptest.NewServer(router)
	defer server.Close()
	serverURL := server.URL
	getReq, err := http.NewRequest(http.MethodHead, serverURL+"/a/b", bytes.NewReader(nil))
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(getReq)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	assert.Equal(t, statusCode, resp.StatusCode)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	assert.NotEqual(t, string(bodyBytes), serverResponse)
	assert.Equal(t, string(bodyBytes), "")
}

func handlePostPutMethod(method string, t *testing.T) {
	requestBody := bytes.Repeat([]byte("a"), 1<<20)
	statusCode := 404
	serverResponse := "serverResponse"

	router := New()
	router.Handle(method, "/a/b", func(c *Context) {
		requestData, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, requestData, requestBody)
		c.Status(statusCode)
		c.String(serverResponse)
	})

	server := httptest.NewServer(router)
	serverURL := server.URL
	defer server.Close()

	getReq, err := http.NewRequest(method, serverURL+"/a/b", bytes.NewReader(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	getReq.Header.Set("Content-Type", "text/plain")
	resp, err := http.DefaultClient.Do(getReq)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, statusCode, resp.StatusCode)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	assert.Equal(t, string(bodyBytes), serverResponse)
}

func handleMethod(method string, t *testing.T) {
	serverResponse := "server response"
	statusCode := 404
	router := New()
	router.Handle(method, "/a/b", func(c *Context) {
		c.Status(statusCode)
		c.String(serverResponse)
	})

	server := httptest.NewServer(router)
	defer server.Close()
	serverURL := server.URL
	getReq, err := http.NewRequest(method, serverURL+"/a/b", bytes.NewReader(nil))
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(getReq)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	assert.Equal(t, statusCode, resp.StatusCode)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	assert.Equal(t, string(bodyBytes), serverResponse)
}

func TestStaticFile(t *testing.T) {
	router := New()
	filePath := "./test/index.html"
	router.StaticFile("/a/b", filePath)

	server := httptest.NewServer(router)
	defer server.Close()
	serverURL := server.URL
	resp, err := http.Get(serverURL + "/a/b")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	sourceFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, bodyBytes, sourceFile)
}

func TestStatic(t *testing.T) {
	router := New()
	dirPath := "./test/"
	fileName := "index.html"
	router.Static("/a/b", dirPath)

	server := httptest.NewServer(router)
	defer server.Close()
	serverURL := server.URL
	resp, err := http.Get(serverURL + "/a/b/" + fileName)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	sourceFile, err := ioutil.ReadFile(dirPath + fileName)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, sourceFile, bodyBytes)
}

func TestNoRoutePrefix(t *testing.T) {
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
}

func TestNoMethod(t *testing.T) {
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
}

func TestPrefix(t *testing.T) {
	router := New()
	serverResponse := "server response"
	statusCode := 404
	v1 := router.Prefix("/v1")
	assert.NotNil(t, v1.engine)
	assert.NotNil(t, v1.router)
	assert.Equal(t, v1.basePath, "/v1")
	v1.Get("/a/b", func(c *Context) {
		c.Status(statusCode)
		c.String(serverResponse)
	})

	server := httptest.NewServer(router)
	defer server.Close()

	serverURL := server.URL
	resp, err := http.Get(serverURL + "/v1/a/b")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	assert.Equal(t, statusCode, resp.StatusCode)

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, serverResponse, string(bodyBytes))
}

func TestPrefixUse(t *testing.T) {
	t.Run("prefix use1", func(t *testing.T) {
		router := New()
		serverResponse := "server response"
		statusCode := 404
		router.Use(func(c *Context) {
			c.SetHeader("fake-header", "fake")
		})

		router.Get("/a/b", func(c *Context) {
			c.Status(statusCode)
			c.String(serverResponse)
		})

		v1 := router.Prefix("/v1")
		v1.Use(func(c *Context) {
			c.SetHeader("version1-header", "version1")
		})

		v1.Get("/a/b", func(c *Context) {
			c.Status(statusCode)
			c.String(serverResponse)
		})

		server := httptest.NewServer(router)
		defer server.Close()

		serverURL := server.URL

		// test router /v1/a/b
		resp, err := http.Get(serverURL + "/v1/a/b")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, statusCode, resp.StatusCode)
		assert.Equal(t, "fake", resp.Header.Get("fake-header"))
		assert.Equal(t, resp.Header.Get("version1-header"), "version1")

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, serverResponse, string(bodyBytes))
		resp.Body.Close()

		// test router /a/b
		resp, err = http.Get(serverURL + "/a/b")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, statusCode, resp.StatusCode)
		assert.Equal(t, "fake", resp.Header.Get("fake-header"))
		assert.Empty(t, resp.Header.Get("version1-header"))

		bodyBytes, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, serverResponse, string(bodyBytes))
		resp.Body.Close()
	})

	t.Run("prefix use2", func(t *testing.T) {
		statusCode := 404
		serverResponse := "server response"
		middleware1 := func(c *Context) {
			assert.Equal(t, c.Header("fake-header"), "fake")
			c.Next()
			assert.Equal(t, c.ResponseWriter.Header().Get("after-request"), "after")
		}
		middleware2 := func(c *Context) {
			c.SetHeader("response-fake-header", "fake")
			c.Next()
			c.String(serverResponse)
		}
		router := New()
		v1 := router.Prefix("/v1")
		v1.Use(middleware1, middleware2)
		v1.Get("/a/b", func(c *Context) {
			c.Status(statusCode)
			c.SetHeader("after-request", "after")
		})

		server := httptest.NewServer(router)
		defer server.Close()
		serverURL := server.URL

		getReq, err := http.NewRequest(http.MethodGet, serverURL+"/v1/a/b", nil)
		if err != nil {
			t.Fatal(err)
		}
		getReq.Header.Set("fake-header", "fake")
		resp, err := http.DefaultClient.Do(getReq)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		assert.Equal(t, statusCode, resp.StatusCode)
		assert.Equal(t, resp.Header.Get("response-fake-header"), "fake")
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, serverResponse, string(bodyBytes))
	})
}

func TestLoadHTMLGlob(t *testing.T) {
	statusCode := 404
	router := New()
	router.LoadHTMLGlob("test/templates/*")
	router.Get("/index.html", func(c *Context) {
		c.Status(statusCode)
		c.HTML("index.tmpl", JSON{
			"title": "Posts",
		})
	})

	server := httptest.NewServer(router)
	defer server.Close()

	serverURL := server.URL
	resp, err := http.Get(serverURL + "/index.html")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	assert.Equal(t, statusCode, resp.StatusCode)

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, strings.Contains(string(bodyBytes), "Posts"))
}

func TestLoadHTMLFiles(t *testing.T) {
	statusCode := 404
	router := New()
	router.LoadHTMLFiles("test/templates/index.tmpl")
	router.Get("/index.html", func(c *Context) {
		c.Status(statusCode)
		c.HTML("index.tmpl", JSON{
			"title": "Posts",
		})
	})

	server := httptest.NewServer(router)
	defer server.Close()

	serverURL := server.URL
	resp, err := http.Get(serverURL + "/index.html")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	assert.Equal(t, statusCode, resp.StatusCode)

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, strings.Contains(string(bodyBytes), "Posts"))
}

func TestUse(t *testing.T) {
	statusCode := 404
	serverResponse := "server response"
	middleware1 := func(c *Context) {
		assert.Equal(t, c.Header("fake-header"), "fake")
		c.Next()
		assert.Equal(t, c.ResponseWriter.Header().Get("after-request"), "after")
	}
	middleware2 := func(c *Context) {
		c.SetHeader("response-fake-header", "fake")
		c.Next()
		c.String(serverResponse)
	}
	router := New()
	router.Use(middleware1, middleware2)
	router.Get("/a/b", func(c *Context) {
		c.Status(statusCode)
		c.SetHeader("after-request", "after")
	})

	server := httptest.NewServer(router)
	defer server.Close()
	serverURL := server.URL

	getReq, err := http.NewRequest(http.MethodGet, serverURL+"/a/b", nil)
	if err != nil {
		t.Fatal(err)
	}
	getReq.Header.Set("fake-header", "fake")
	resp, err := http.DefaultClient.Do(getReq)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	assert.Equal(t, statusCode, resp.StatusCode)
	assert.Equal(t, resp.Header.Get("response-fake-header"), "fake")
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, serverResponse, string(bodyBytes))
}

type testHandler struct {
	HandlerFunc func(*Context)
}

func (handler *testHandler) Handle(c *Context) {
	handler.HandlerFunc(c)
}

func TestUseHandler(t *testing.T) {
	statusCode := 404
	serverResponse := "server response"

	handler1 := &testHandler{
		HandlerFunc: func(c *Context) {
			assert.Equal(t, c.Header("fake-header"), "fake")
			c.Next()
			assert.Equal(t, c.ResponseWriter.Header().Get("after-request"), "after")
		},
	}

	handler2 := &testHandler{
		HandlerFunc: func(c *Context) {
			c.SetHeader("response-fake-header", "fake")
			c.Next()
			c.String(serverResponse)
		},
	}

	router := New()
	router.UseHandler(handler1, handler2)
	router.Get("/a/b", func(c *Context) {
		c.Status(statusCode)
		c.SetHeader("after-request", "after")
	})

	server := httptest.NewServer(router)
	defer server.Close()
	serverURL := server.URL

	getReq, err := http.NewRequest(http.MethodGet, serverURL+"/a/b", nil)
	if err != nil {
		t.Fatal(err)
	}
	getReq.Header.Set("fake-header", "fake")
	resp, err := http.DefaultClient.Do(getReq)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	assert.Equal(t, statusCode, resp.StatusCode)
	assert.Equal(t, resp.Header.Get("response-fake-header"), "fake")
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, serverResponse, string(bodyBytes))
}
