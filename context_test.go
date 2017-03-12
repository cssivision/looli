package looli

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestQuery(t *testing.T) {
	statusCode := 200
	serverResponse := "server response"
	router := New()
	router.Get("/path", func(c *Context) {
		assert.Empty(t, c.Query(""))
		assert.Equal(t, c.Query("name"), "cssivision")
		assert.Equal(t, c.Query("age"), "23")
		assert.Equal(t, c.Query("bar"), "イモト")
		assert.Empty(t, c.Query("other"))
		c.Status(statusCode)
		c.String(serverResponse)
	})

	server := httptest.NewServer(router)
	defer server.Close()
	serverURL := server.URL
	resp, err := http.Get(serverURL + "/path?&name=cssivision&age=23&bar=イモト")
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

func TestDefaultQuery(t *testing.T) {
	statusCode := 200
	serverResponse := "server response"
	router := New()
	router.Get("/path", func(c *Context) {
		assert.Empty(t, c.DefaultQuery("", ""))
		assert.Equal(t, c.DefaultQuery("name", "biz"), "cssivision")
		assert.Equal(t, c.DefaultQuery("age", "24"), "23")
		assert.Equal(t, c.DefaultQuery("other", "other value"), "other value")
		c.Status(statusCode)
		c.String(serverResponse)
	})

	server := httptest.NewServer(router)
	defer server.Close()
	serverURL := server.URL
	resp, err := http.Get(serverURL + "/path?&name=cssivision&age=23")
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

func TestPostForm(t *testing.T) {
	statusCode := 200
	serverResponse := "server response"
	router := New()
	router.Post("/path", func(c *Context) {
		assert.Equal(t, c.PostForm("foo"), "bar")
		assert.Equal(t, c.PostForm("page"), "11")
		assert.Empty(t, c.PostForm("both"))
		assert.Empty(t, c.PostForm("other"))
		c.Status(statusCode)
		c.String(serverResponse)
	})

	server := httptest.NewServer(router)
	defer server.Close()
	serverURL := server.URL

	body := bytes.NewReader([]byte("foo=bar&page=11&both=&foo=second"))
	resp, err := http.Post(serverURL+"/path", MIMEPOSTForm, body)
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

func TestDefaultPostForm(t *testing.T) {
	statusCode := 200
	serverResponse := "server response"
	router := New()
	router.Post("/path", func(c *Context) {
		assert.Equal(t, c.DefaultPostForm("foo", "hh"), "bar")
		assert.Equal(t, c.DefaultPostForm("page", "12"), "11")
		assert.Equal(t, c.DefaultPostForm("both", "other"), "other")
		assert.Empty(t, c.DefaultPostForm("other", ""))
		c.Status(statusCode)
		c.String(serverResponse)
	})

	server := httptest.NewServer(router)
	defer server.Close()
	serverURL := server.URL

	body := bytes.NewReader([]byte("foo=bar&page=11&both=&foo=second"))
	resp, err := http.Post(serverURL+"/path", MIMEPOSTForm, body)
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

func TestPostFormMultipart(t *testing.T) {
	boundary := "--testboundary"
	body := new(bytes.Buffer)
	mw := multipart.NewWriter(body)
	must(mw.SetBoundary(boundary))
	must(mw.WriteField("foo", "bar"))
	must(mw.WriteField("bar", "10"))
	must(mw.WriteField("array", "first"))
	must(mw.WriteField("id", "12"))
	mw.Close()

	statusCode := 200
	serverResponse := "server response"
	router := New()
	router.Post("/a/b", func(c *Context) {
		assert.Equal(t, c.PostForm("bar"), "10")
		assert.Equal(t, c.PostForm("foo"), "bar")
		assert.Equal(t, c.PostForm("array"), "first")
		assert.Equal(t, c.PostForm("id"), "12")
		c.Status(statusCode)
		c.String(serverResponse)
	})

	server := httptest.NewServer(router)
	defer server.Close()

	serverURL := server.URL
	getReq, err := http.NewRequest(http.MethodPost, serverURL+"/a/b", body)
	if err != nil {
		t.Fatal(err)
	}
	getReq.Header.Set("Content-Type", MIMEMultipartPOSTForm+"; boundary="+boundary)
	resp, err := http.DefaultClient.Do(getReq)
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

func must(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func TestStatus(t *testing.T) {
	statusCode := 404
	serverResponse := "server response"
	router := New()
	router.Get("/", func(c *Context) {
		c.Status(statusCode)
		c.String(serverResponse)
	})

	server := httptest.NewServer(router)
	defer server.Close()
	serverURL := server.URL
	resp, err := http.Get(serverURL)
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

func TestRedirect(t *testing.T) {
	statusCode := 404
	serverResponse := "server response"
	router := New()
	router.Get("/redirect", func(c *Context) {
		c.Status(statusCode)
		c.String(serverResponse)
	})
	router.Get("/", func(c *Context) {
		c.Redirect("/redirect")
	})

	server := httptest.NewServer(router)
	defer server.Close()
	serverURL := server.URL
	resp, err := http.Get(serverURL)
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

func TestServeFile(t *testing.T) {
	router := New()
	router.Get("/a/b", func(c *Context) {
		c.ServeFile("test/index.html")
	})

	server := httptest.NewServer(router)
	defer server.Close()

	serverURL := server.URL
	resp, err := http.Get(serverURL + "/a/b")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadFile("test/index.html")
	if err != nil {
		t.Fatal(err)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, content, bodyBytes)
}

func TestHeader(t *testing.T) {
	statusCode := 404
	serverResponse := "server response"
	router := New()
	router.Get("/a/b", func(c *Context) {
		assert.Equal(t, c.Header("fake-header"), "fake")
		c.Status(statusCode)
		c.String(serverResponse)
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
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, serverResponse, string(bodyBytes))
}

func TestSetHeader(t *testing.T) {
	statusCode := 404
	serverResponse := "server response"
	router := New()
	router.Get("/a/b", func(c *Context) {
		c.SetHeader("fake-header", "fake")
		c.Status(statusCode)
		c.String(serverResponse)
	})

	server := httptest.NewServer(router)
	defer server.Close()

	serverURL := server.URL
	resp, err := http.Get(serverURL + "/a/b")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	assert.Equal(t, resp.Header.Get("fake-header"), "fake")
	assert.Equal(t, statusCode, resp.StatusCode)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, serverResponse, string(bodyBytes))
}

func TestCookie(t *testing.T) {
	statusCode := 404
	serverResponse := "server response"
	router := New()
	router.Get("/a/b", func(c *Context) {
		val, err := c.Cookie("fake-cookie")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, val, "fake")
		c.Status(statusCode)
		c.String(serverResponse)
	})

	server := httptest.NewServer(router)
	defer server.Close()

	serverURL := server.URL
	getReq, err := http.NewRequest(http.MethodGet, serverURL+"/a/b", nil)
	if err != nil {
		t.Fatal(err)
	}

	getReq.AddCookie(&http.Cookie{
		Name:  "fake-cookie",
		Value: "fake",
	})
	resp, err := http.DefaultClient.Do(getReq)
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

func TestSetCookie(t *testing.T) {
	statusCode := 404
	serverResponse := "server response"
	router := New()
	router.Get("/a/b", func(c *Context) {
		c.SetCookie(&http.Cookie{
			Name:  "fake-cookie",
			Value: "fake",
		})
		c.Status(statusCode)
		c.String(serverResponse)
	})

	server := httptest.NewServer(router)
	defer server.Close()

	serverURL := server.URL
	resp, err := http.Get(serverURL + "/a/b")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	assert.Equal(t, "fake", resp.Cookies()[0].Value)
	assert.Equal(t, statusCode, resp.StatusCode)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, serverResponse, string(bodyBytes))
}

func TestAbort(t *testing.T) {
	statusCode := 404
	serverResponse := "server response"
	middleware1 := func(c *Context) {
		c.SetHeader("fake-header4", "fake4")
		c.Next()
	}

	middleware2 := func(c *Context) {
		c.SetHeader("fake-header3", "fake3")
		c.Abort()
		assert.True(t, c.IsAborted())
		assert.Equal(t, c.current, abortIndex)
		c.Status(statusCode)
		c.String(serverResponse)
	}

	middleware3 := func(c *Context) {
		c.SetHeader("fake-header2", "fake2")
	}

	router := New()
	router.Use(middleware1, middleware2, middleware3)
	router.Get("/a/b", func(c *Context) {
		c.SetHeader("fake-header1", "fake1")
	})

	server := httptest.NewServer(router)
	defer server.Close()

	serverURL := server.URL
	resp, err := http.Get(serverURL + "/a/b")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, statusCode, resp.StatusCode)
	assert.Empty(t, resp.Header.Get("fake-header1"))
	assert.Empty(t, resp.Header.Get("fake-header2"))
	assert.Equal(t, "fake3", resp.Header.Get("fake-header3"))
	assert.Equal(t, "fake4", resp.Header.Get("fake-header4"))
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, serverResponse, string(bodyBytes))
}

func TestAbortWithStatus(t *testing.T) {
	statusCode := 404
	serverResponse := "server response"
	middleware1 := func(c *Context) {
		c.AbortWithStatus(statusCode)
		assert.True(t, c.IsAborted())
		assert.Equal(t, c.current, abortIndex)
	}
	middleware2 := func(c *Context) {
		c.String(serverResponse)
	}

	router := New()
	router.Use(middleware1, middleware2)
	router.Get("/a/b", func(c *Context) {})

	server := httptest.NewServer(router)
	defer server.Close()

	serverURL := server.URL
	resp, err := http.Get(serverURL + "/a/b")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, statusCode, resp.StatusCode)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	assert.Empty(t, bodyBytes)
}

func TestParam(t *testing.T) {
	statusCode := 404
	serverResponse := "server response"
	router := New()
	router.Get("/a/:name", func(c *Context) {
		assert.Equal(t, "cssivision", c.Param("name"))
		assert.Empty(t, c.Param("other"))
		c.Status(statusCode)
		c.String(serverResponse)
	})

	router.Get("/b/*filepath", func(c *Context) {
		assert.Equal(t, "c/cssivision", c.Param("filepath"))
		c.Status(statusCode)
		c.String(serverResponse)
	})

	server := httptest.NewServer(router)
	defer server.Close()
	serverURL := server.URL

	resp, err := http.Get(serverURL + "/a/cssivision")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, statusCode, resp.StatusCode)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, serverResponse, string(bodyBytes))
	resp.Body.Close()

	resp, err = http.Get(serverURL + "/b/c/cssivision")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, statusCode, resp.StatusCode)
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, serverResponse, string(bodyBytes))
	resp.Body.Close()
}

func TestClientIP(t *testing.T) {
	t.Run("X-Real-Ip", func(t *testing.T) {
		statusCode := 404
		serverResponse := "server response"
		realIP := "looli.xyz"
		router := New()
		router.ForwardedByClientIP = true

		router.Get("/a", func(c *Context) {
			assert.Equal(t, realIP, c.ClientIP())
			c.Status(statusCode)
			c.String(serverResponse)
		})

		server := httptest.NewServer(router)
		defer server.Close()

		serverURL := server.URL
		getReq, err := http.NewRequest(http.MethodGet, serverURL+"/a", nil)
		assert.Nil(t, err)
		getReq.Header.Set("X-Real-Ip", realIP)
		resp, err := http.DefaultClient.Do(getReq)
		assert.Nil(t, err)
		defer resp.Body.Close()

		assert.Equal(t, statusCode, resp.StatusCode)
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		assert.Nil(t, err)
		assert.Equal(t, serverResponse, string(bodyBytes))
	})

	t.Run("X-Forwarded-For", func(t *testing.T) {
		statusCode := 404
		serverResponse := "server response"
		clientIP := "looli.xyz, looli.com"
		router := New()
		router.ForwardedByClientIP = true

		router.Get("/a", func(c *Context) {
			assert.Equal(t, "looli.xyz", c.ClientIP())
			assert.Empty(t, c.Header("X-Real-Ip"))
			c.Status(statusCode)
			c.String(serverResponse)
		})

		server := httptest.NewServer(router)
		defer server.Close()

		serverURL := server.URL
		getReq, err := http.NewRequest(http.MethodGet, serverURL+"/a", nil)
		assert.Nil(t, err)
		getReq.Header.Set("X-Forwarded-For", clientIP)
		resp, err := http.DefaultClient.Do(getReq)
		assert.Nil(t, err)
		defer resp.Body.Close()

		assert.Equal(t, statusCode, resp.StatusCode)
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		assert.Nil(t, err)
		assert.Equal(t, serverResponse, string(bodyBytes))
	})
}

func TestContentType(t *testing.T) {
	statusCode := 404
	serverResponse := "server response"
	router := New()
	router.Post("/a/b", func(c *Context) {
		assert.Equal(t, "text/plain", c.ContentType())
		c.Status(statusCode)
		c.String(serverResponse)
	})

	server := httptest.NewServer(router)
	defer server.Close()

	serverURL := server.URL
	getReq, err := http.NewRequest(http.MethodPost, serverURL+"/a/b", nil)
	getReq.Header.Set("Content-Type", "text/plain")
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
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, serverResponse, string(bodyBytes))
}

func TestString(t *testing.T) {
	statusCode := 404
	serverResponse := "server response"
	router := New()
	router.Get("/a/b", func(c *Context) {
		c.Status(statusCode)
		c.String(serverResponse)
	})

	server := httptest.NewServer(router)
	defer server.Close()

	serverURL := server.URL
	resp, err := http.Get(serverURL + "/a/b")
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

func TestJSON(t *testing.T) {
	type Info struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	statusCode := 404
	router := New()
	router.Get("/a/b", func(c *Context) {
		c.Status(statusCode)
		c.JSON(JSON{
			"name": "cssivision",
			"age":  21,
		})
	})

	server := httptest.NewServer(router)
	defer server.Close()

	serverURL := server.URL
	resp, err := http.Get(serverURL + "/a/b")
	assert.Nil(t, err)
	defer resp.Body.Close()
	assert.Equal(t, statusCode, resp.StatusCode)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	data := new(Info)
	err = json.Unmarshal(bodyBytes, data)
	assert.Nil(t, err)
	assert.Equal(t, "cssivision", data.Name)
	assert.Equal(t, 21, data.Age)
}

func TestHTML(t *testing.T) {
	t.Run("normal render", func(t *testing.T) {
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
		assert.Nil(t, err)
		defer resp.Body.Close()
		assert.Equal(t, statusCode, resp.StatusCode)
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		assert.Nil(t, err)
		assert.True(t, strings.Contains(string(bodyBytes), "Posts"))
	})

	t.Run("render panic", func(t *testing.T) {
		statusCode := 404
		router := New()
		router.LoadHTMLGlob("test/templates/*")
		router.Get("/index.html", func(c *Context) {
			c.Status(statusCode)
			assert.Panics(t, func() {
				c.HTML("index.tmp", JSON{
					"title": "Posts",
				})
			})
		})

		server := httptest.NewServer(router)
		defer server.Close()

		serverURL := server.URL
		resp, err := http.Get(serverURL + "/index.html")
		assert.Nil(t, err)
		defer resp.Body.Close()
		assert.Equal(t, statusCode, resp.StatusCode)
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		assert.Nil(t, err)
		assert.False(t, strings.Contains(string(bodyBytes), "Posts"))
	})
}
